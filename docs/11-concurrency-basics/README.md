# 第11章：并发基础

## 学习目标
- 创建 goroutine 并理解调度
- 使用 WaitGroup 同步结束
- 识别数据竞争并避免共享可变状态
- 知道何时使用锁与何时使用通道

## 章节提纲
- go 关键字启动并发任务
- GOMAXPROCS 与协作式调度
- sync.WaitGroup 的基本同步
- 竞争条件示例与 go run -race
- 锁的选择：Mutex/RWMutex 简介

## goroutine 与调度
- `go f()` 创建 goroutine，立即返回，不保证执行顺序；主 goroutine 退出会导致进程结束。
- Go 运行时使用 M:N 调度：多个 goroutine 复用少量内核线程；`GOMAXPROCS` 控制同时运行的 P（逻辑 CPU）数量，默认等于 CPU 核数。
- 调度是协作式的，goroutine 在系统调用、阻塞、channel/锁等待、函数调用栈切换点等位置挂起；不要依赖“睡一会让它先跑”。

示例：
```go
go fmt.Println("hello from goroutine")
fmt.Println("main")
time.Sleep(10 * time.Millisecond) // 仅为示例确保 goroutine 有机会运行
```

## WaitGroup 基本同步
- 适合等待一组 goroutine 结束：`Add` 计数，goroutine 内 `Done`，主线程 `Wait`。
- 使用原则：`Add` 在启动 goroutine 之前调用；不要拷贝 WaitGroup 值；每个 `Add(1)` 对应一个 `Done()`。

示例：
```go
var wg sync.WaitGroup
wg.Add(3)
for i := 0; i < 3; i++ {
	go func(id int) {
		defer wg.Done()
		work(id)
	}(i)
}
wg.Wait()
```

## 数据竞争与 race 检测
- 数据竞争定义：多个 goroutine 同时读写同一内存，且至少有一个写，且缺少同步。
- 典型陷阱：共享变量自增、共享 map 写、循环变量捕获。
- 用 `go test -race` 或 `go run -race main.go` 检测；发现问题后用锁、channel 或复制数据消除共享。

示例（有竞争，不要模仿）：
```go
var counter int
go func() { counter++ }()
go func() { counter++ }()
fmt.Println(counter)
```
修复：用 `sync.Mutex` 保护，或使用 channel 串行化增量。

## 锁的选择与用法
- `sync.Mutex`：独占锁，最常用。`Lock`/`Unlock` 成对；`defer Unlock()` 便于异常路径。
- `sync.RWMutex`：读多写少时，读可并发、写仍独占；滥用会增加开销和死锁复杂度。
- `sync.Mutex` 不可复制，通常嵌入结构体指针使用。

示例：
```go
type Counter struct {
	mu sync.Mutex
	n  int
}
func (c *Counter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.n++
}
func (c *Counter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.n
}
```

## 何时用 channel，何时用锁（概要）
- channel 强调“通过通信共享内存”：在管道式传递数据、事件通知、背压时使用；更多语义与模式放在第 12 章。
- 锁保护共享可变状态：需要原地修改的结构、缓存、计数器时使用。
- 判断标准：数据是否天然要移动？移动/复制容易则 channel；必须原地保护则锁。

## 常见模式与技巧（基础版）
- 超时/取消：`context.WithTimeout` 或 `select` + `time.After`；长操作要响应 `ctx.Done()`。
- 控制并发度：worker pool（示例见本章代码）、带缓冲 channel 充当信号量。
- 循环变量捕获：循环内 `v := v` 再传给 goroutine。
- 清理：goroutine 内 `defer` 关闭资源；主流程用 WaitGroup 等待，避免泄漏。

## 补充：第三方 goroutine 池（如 ants）
- 当需要大量短任务且频繁创建/销毁 goroutine 时，可使用 goroutine 池减少调度和 GC 压力。常用库：`github.com/panjf2000/ants`。
- 核心思路：预设最大并发，任务提交时由池复用现有 goroutine；可配置池大小、超时、Panic 处理等。
- 使用示例：
```go
pool, _ := ants.NewPool(10)            // 最大并发 10
defer pool.Release()

for i := 0; i < 100; i++ {
	v := i
	_ = pool.Submit(func() {
		fmt.Println("work", v)
	})
}
pool.Wait() // 等待任务完成（需使用 WithPreAlloc(false) 默认开启阻塞等待）
```
- 何时考虑：在高 QPS、短任务、对延迟敏感的场景，或需限制后台 goroutine 数量时；简单场景直接 `go func()` 更易读。

## 继续：通道与模式（详见第 12 章）
- 想系统学习 channel 语义、select 用法、pipeline、fan-out/fan-in、背压与关闭约定，请阅读第 12 章。
- 可先完成本章小作业，再在第 12 章中实现基于 channel 的版本，对比锁与通信的差异。
## 小作业
1) 写一个 worker pool：限定并发 N，处理任务切片，汇总结果；用 WaitGroup 同步结束。  
2) 写一个带超时的请求函数：使用 `context.WithTimeout`，goroutine 模拟处理，超时返回错误。  
3) 改写一个有数据竞争的计数器，用锁或 channel 修复，并用 `go test -race` 验证。
