# 第12章：通道与并发模式

## 学习目标
- 掌握通道的发送/接收/关闭语义
- 利用缓冲通道与 select 构建模式
- 设计 pipeline 与 fan-out/fan-in
- 用关闭通道传播完成信号

## 章节提纲
- 无缓冲 vs 有缓冲通道的特性
- select 用于超时、合并与取消
- 流水线、fan-out/fan-in 模式示例
- 生产者/消费者与背压控制
- 通道关闭约定与避免双关闭

> 前置：默认已掌握第 11 章的 goroutine、WaitGroup、锁与 race 检测。本章专注 channel/ select 与并发模式。

## 何时使用 go、何时使用 channel
- 用 `go`：当任务可以异步执行且不依赖立即结果；常见于并行 I/O、后台通知、独立清理任务。启动后要有结束条件（ctx/done）避免 goroutine 泄漏。
- 用 channel：在 goroutine 间**传递数据/事件**或**组合并发流程**时；如果只是保护共享状态更适合锁。
- “只 go 不沟通”会导致失控：每个 goroutine 都需要退出通道/ctx、错误返回或 WaitGroup 等同步手段。
- “只用 channel 不合适”场景：需要原地修改共享结构、缓存、计数等；使用锁更简单直接。

## 无缓冲 vs 有缓冲
- 无缓冲：发送与接收同步握手；可用于任务交接、限速；发送方阻塞直到有人接收。
- 有缓冲：`make(chan T, n)`；缓冲未满时发送不阻塞，缓冲为空时接收阻塞；用于平滑突发、实现 semaphore。
- 不要用 `len(ch)` 作为逻辑条件（存在竞态）；用 select 或额外信号通道表达状态。

## select 的常见用法
- 超时：`case <-time.After(d):`，为每次 select 创建新计时器；高频时用 `time.NewTimer` 复用。
- 取消：`case <-ctx.Done():` 及时退出；上游关闭通道/取消 context 传播结束。
- 合并多个来源：`select` 监听多个输入通道，将数据聚合到一个输出。
- 非阻塞操作：加入 `default` 分支，用于尝试发送/接收失败则跳过。

示例：超时等待
```go
select {
case v := <-ch:
	fmt.Println("got", v)
case <-time.After(500 * time.Millisecond):
	return errors.New("timeout")
}
```

## 流水线（pipeline）
- 将任务拆为若干阶段，每阶段接收一个通道并输出到下一个；便于并行与解耦。
```go
stage1 := func(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for v := range in {
			out <- v * 2
		}
	}()
	return out
}
stage2 := func(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for v := range in {
			out <- v + 1
		}
	}()
	return out
}

in := make(chan int, 3)
in <- 1; in <- 2; in <- 3
close(in)
for v := range stage2(stage1(in)) {
	fmt.Println(v) // 3,5,7
}
```

## fan-out / fan-in
- fan-out：多个 worker 从同一输入通道消费，提高并行度。
- fan-in：合并多个输出通道到一个通道，常配合 WaitGroup/计数器关闭。

示例（简化版）：
```go
worker := func(id int, jobs <-chan int, out chan<- int) {
	for j := range jobs {
		out <- j * 2
	}
}

jobs := make(chan int)
out := make(chan int)
for i := 0; i < 3; i++ {
	go worker(i, jobs, out)
}

go func() {
	defer close(jobs)
	for i := 0; i < 5; i++ {
		jobs <- i
	}
}()

go func() {
	defer close(out)
	// 等待所有 worker 退出通常用 WaitGroup；此处简化忽略。
}()

for v := range out {
	fmt.Println(v)
}
```

## 生产者/消费者与背压
- 背压：下游处理慢时让上游阻塞或丢弃；缓冲 channel 是天然的背压工具。
- 尽量明确“谁关闭通道”：生产者关闭输出通道；消费者不关闭它未创建的通道，避免双关闭。
- 当需要退出时，首选 context 或专门的 `done` 通道广播结束。

## 通道关闭约定
- 发送端关闭：表示“不再有新数据”，接收端仍可读到剩余缓冲。
- 判断关闭：`v, ok := <-ch`；`ok==false` 表示通道已关闭且无数据。
- 不要向已关闭的通道发送，会 panic；也不要重复关闭。

## 调试建议
- 小心 goroutine 泄漏：在 select 中始终处理 `ctx.Done()` 或 done 通道。
- race 检测：`go test -race` 能捕获对通道/共享变量的并发误用。

## 实践作业
1) Pipeline：实现两阶段流水线（如 `x*2` 再 `x+1`），用 channel 串联并写测试验证输出顺序。  
2) fan-out/fan-in：创建 N 个 worker 处理 jobs（如平方），使用 WaitGroup 收拢并关闭输出通道，测试结果数量与内容。  
3) 背压与超时：实现一个带超时的发送函数（`select` + `time.After` 或 `ctx.Done()`），当 channel 满或取消时返回错误；为其编写测试验证超时/成功路径。
