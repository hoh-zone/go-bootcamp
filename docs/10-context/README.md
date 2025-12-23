# 第10章：Context 与取消

## 学习目标
- 使用 context 传递取消/超时
- 在 I/O 与 goroutine 中监听 Done
- 把 context 融入接口与 HTTP 处理
- 避免在 context 中存放大对象

## 章节提纲
- Background/TODO 与派生上下文
- WithCancel/Timeout/Deadline 用法
- select ctx.Done 处理退出与泄漏防护
- context 在库/接口设计中的位置
- 不在 context 中存放大数据的原则

## 基本类型与派生
- 根 context：`context.Background()`（服务器/CLI 入口使用），`context.TODO()`（暂时占位，尽快替换）。
- 派生：`WithCancel`（手动取消）、`WithTimeout`/`WithDeadline`（自动取消）、`WithValue`（谨慎传少量元数据）。
- 取消是可传递的：子 context 取消不会影响父，父取消会立即通知所有子。

示例：
```go
ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
defer cancel()
doWork(ctx)
```

## “结构体”是怎么组织的：一条 Context 链
```go
ctx := context.Background()
ctx = context.WithValue(ctx, k, v)
ctx, cancel := context.WithCancel(ctx)
ctx, cancel2 := context.WithTimeout(ctx, time.Second)
```
在内部就是一条链，大致像这样：

```
timerCtx
  └─ cancelCtx
       └─ valueCtx
            └─ emptyCtx(background)
```
每个 wrapper struct 都保存一个 parent Context，并且除了自己负责的那一项，其它方法都转发给 parent（delegation）。

## 在 goroutine 与 I/O 中监听 Done
- 长时间运行的 goroutine、阻塞 I/O、循环都要检查 `ctx.Done()`，退出时释放资源。
- 使用 `select` 同时等待业务通道与 `ctx.Done()`，避免 goroutine 泄漏。

示例：
```go
func worker(ctx context.Context, jobs <-chan Job) {
	for {
		select {
		case j, ok := <-jobs:
			if !ok { //从通道接收时该通道已被关闭（且缓冲被读空）
				return
			}
			process(j)
		case <-ctx.Done():
			return
		}
	}
}
```

## WithTimeout/Deadline 的实践
- 为外部依赖（HTTP、数据库、队列等）设定超时，防止无限阻塞。
- 注意：`time.After` 每次会分配计时器，高频情况下用 `time.NewTimer` 并复用。
- `context.DeadlineExceeded` 与 `context.Canceled` 是常见错误检查。

## 在接口与 HTTP 处理中使用 context
- 将 `ctx` 作为第一个参数传递：`func (s *Service) Fetch(ctx context.Context, id string) error`，方便调用栈统一取消。
- HTTP：`r.Context()` 获取请求的 ctx；下游调用都用该 ctx；在 handler 里尊重取消，避免继续写响应。
- gRPC、database/sql 等库自动支持 ctx；确保传递上游 ctx 而非新建无关联的 Background。

## WithValue 的使用边界
- 仅存放与请求范围相关的元数据（如 request-id、用户身份）；避免大对象和可变数据。
- key 应使用自定义类型避免冲突：`type ctxKey string`。
- 不要用 ctx 传递必需的业务参数，保持函数显式参数。

## 常见陷阱
- 忘记调用 `cancel()`：导致计时器和资源泄漏；用 `defer cancel()` 靠近创建处。
- 在函数内重新 `context.Background()`：会丢失上游取消信号，应该接受并传递入参 ctx。
- 在并发 map/全局变量中存 ctx：ctx 设计为请求级别，完成后应丢弃。
- 重复关闭结果通道时未监听 ctx：结合 WaitGroup/关闭方约定，避免 panic。

## 调试与测试
- 检查取消路径：构造超时/取消测试，验证 goroutine 提前退出。
- 在模拟依赖时，用 `select` + `ctx.Done()` 让 fake 也可被取消。
- race 检测仍有价值：确保取消后不再访问已释放资源。
