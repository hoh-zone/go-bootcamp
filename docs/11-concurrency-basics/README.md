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

## 操作系统背景知识

- 进程（Process）：资源拥有者，拥有独立虚拟地址空间、文件句柄等；进程间内存隔离，通信需 IPC。
- 线程（Thread）：CPU 调度的最小单位，运行在进程内，线程间共享进程资源；切换代价小于进程。
- 为什么要多线程：利用多核并行、隔离阻塞（IO、系统调用）、避免单线程长任务卡住整体；但共享内存带来同步成本和风险。
- 锁（Lock）：互斥同步原语，保证同一时刻只有一个执行流进入临界区；常见有互斥锁、读写锁。
- 死锁（Deadlock）：多个执行流相互等待永不释放，典型由“占有且等待、不可抢占、循环等待”组合引起；避免手段包括固定锁顺序、超时/尝试锁、减少持锁时间。
- 信号量（Semaphore）：计数型同步原语，控制并发许可数量；二元信号量可等价互斥锁，计数信号量可用于连接池、限流。

### 用户级线程 vs 内核级线程（简表）
| 特性 | 用户级线程（green thread） | 内核级线程 |
| --- | --- | --- |
| 创建/切换成本 | 低，纯用户态，不触发系统调用 | 较高，需内核参与 |
| 调度者 | 运行时/用户态库 | 操作系统内核 |
| 阻塞系统调用影响 | 阻塞整个内核线程（需运行时规避或切换） | 仅阻塞当前线程，其他线程可继续 |
| 并行性 | 单内核线程下仅并发；要并行需绑多个内核线程 | 天然可在多核并行 |
| 可见性/诊断 | 对 OS 不可见，系统工具支持有限 | OS 可见，诊断工具丰富 |
| 典型例子 | Go goroutine（M:N）、Java 早期 green thread | POSIX 线程、现代 Java 线程 |



## GMP 模型概览
- G（goroutine）、M（内核线程）、P（逻辑处理器）三者组合完成调度：M 必须持有 P 才能运行 G，`GOMAXPROCS` 即 P 的数量上限。
- 每个 P 有本地运行队列，队列空时会从其他 P “偷” 一半任务，或从全局队列取任务；新建 goroutine 优先入当前 P 的本地队列，减少竞争。
- 阻塞处理：M 在系统调用/阻塞时释放手中的 P，P 会绑定到新的 M 继续执行其他 G；当阻塞返回时，G 会重新排队等待调度。
- 网络 IO 由 netpoll 管理：M 在 epoll/kqueue 等事件就绪时唤醒对应 G，避免忙等。

### GMP 运行流程（细化）
- 创建 goroutine：`go f()` 将 G 放入当前 P 的本地队列；队列满则部分推送到全局队列。
- 取任务运行：持有 P 的 M 先从 P 的本地队列取 G，空则从全局队列或其他 P 偷一半。
- 系统调用阻塞：M 进入内核阻塞，先交出 P；调度器唤起新的 M+P 继续跑其他 G。阻塞返回后原 G 重新入队等待。
- 垃圾回收：STW 期间可能调整 P 的可运行状态；并发标记阶段由后台 G + mutator 协作。
- netpoll：网络 IO 注册到内核事件；事件就绪时把对应 G 放入全局可运行队列，空闲 M 会被唤醒。

### GMP 示意
```
          +---------------------------+
          |         Global Queue      |
          +-------------+-------------+
                        |
     +---------+   +---------+   +---------+
     |   P0    |   |   P1    |   |   P2    | ... (P count = GOMAXPROCS)
     | local Q |   | local Q |   | local Q |
     +----+----+   +----+----+   +----+----+
          |             |             |
     +----v----+   +----v----+   +----v----+
     |   M0    |   |   M1    |   |   M2    | ... (M = OS threads)
     +----+----+   +----+----+   +----+----+
          |             |             |
       running G     running G     running G

阻塞时：Mx 进入 syscall，先交出 Px，调度器唤醒/创建新的 My 与 Px 继续运行；Mx 返回后其 G 重回队列。
偷取：某 P 队列空时，从其他 P 的队列尾偷一半，或从 Global Q 取。
```

参考: https://go.cyub.vip/gmp/gmp-model/

### 容器环境下的 GOMAXPROCS
- 容器 CPU 配额/亲和性通常小于宿主机核数，默认 GOMAXPROCS=宿主机核数会创建过多 P，浪费调度并拉高抖动。
- 症状：配额 1 核却有 8 个 P，goroutine 抢占排队、系统调用放大、尾延迟升高。
- 方案：启动时依据 cgroup 配额设置 `runtime.GOMAXPROCS(n)`，或使用自动探测库。
- 推荐库：`github.com/uber-go/automaxprocs` 会在 `init()` 里根据 cgroup/CPU quota 自动调用 `runtime.GOMAXPROCS`，并打印调整结果。
- 用法：
```go
import _ "go.uber.org/automaxprocs"
```
- 诊断：可临时设置 `GODEBUG=cpu.all=1` 查看内核可见 CPU，或打印 `runtime.GOMAXPROCS(0)` 确认最终值。

## goroutine 与调度
- `go f()` 创建 goroutine，立即返回，不保证执行顺序；主 goroutine 退出会导致进程结束。


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

> Go proverb: “Do not communicate by sharing memory; instead, share memory by communicating.” 
> 不要以共享内存的方式来通信，相反，要通过通信来共享内存

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
