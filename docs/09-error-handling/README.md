# 第9章：错误处理与恢复

## 学习目标
- 使用标准库 errors 创建与包装错误
- 区分 sentinel error 与自定义类型
- 掌握 errors.Is/As 解包链路
- 了解 panic/recover 的边界

## 章节提纲
- errors.New 与 fmt.Errorf("%w")
- 自定义错误类型携带上下文
- errors.Is/As 处理错误链
- 日志与错误的分层职责
- panic/recover 的使用边界

## 创建与包装错误
- 基础：`errors.New("msg")`，`fmt.Errorf("extra: %w", err)` 进行包装。
- 包装链路保留根因，方便上层判断与日志。

示例：
```go
var ErrNotFound = errors.New("not found")

func load(id int) (string, error) {
	if id == 0 {
		return "", fmt.Errorf("load %d: %w", id, ErrNotFound)
	}
	return "ok", nil
}
```

## 自定义错误类型
- 当需要携带结构化上下文时，定义错误类型并实现 `Error() string`。
- 通过 `errors.As` 解包到自定义类型，方便读取字段。

示例：
```go
type ParseError struct {
	Line int
	Msg  string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("line %d: %s", e.Line, e.Msg)
}

func parse(line int, s string) error {
	if s == "" {
		return &ParseError{Line: line, Msg: "empty"}
	}
	return nil
}
```

## errors.Is 与 errors.As
- `errors.Is(err, target)` 用于匹配链路上的哨兵错误或实现 `Is` 方法的类型。
- `errors.As(err, &target)` 将链路中第一个可赋值给 `target` 的错误赋给它。

示例：
```go
if err := load(0); err != nil {
	if errors.Is(err, ErrNotFound) {
		fmt.Println("not found")
	}
	var perr *ParseError
	if errors.As(err, &perr) {
		fmt.Println("parse error at", perr.Line)
	}
}
```

## 日志与错误职责
- 函数返回错误时，不要在每一层都重复日志；在“边界”处统一记录（如 handler/worker）。
- 返回的错误应包含上下文，日志再补充请求范围信息（trace id 等）。

## panic/recover 边界
- `panic` 用于不可恢复的编程错误；业务错误应通过返回 `error` 处理。
- `recover` 只能在延迟函数中生效；常见于服务的最外层兜底或 goroutine 入口处。
- 使用 `defer` + `recover` 时记得记录堆栈：`debug.Stack()`。

示例：
```go
func safeRun(fn func()) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic: %v\n%s", r, debug.Stack())
		}
	}()
	fn()
}
```

## 小练习
1) 定义一个 `ValidationError`，包含字段名 `Field` 与原因 `Reason`，实现 `Error()`。
2) 写函数 `WrapLoad(id int) error`，底层调用 `load` 并用 `fmt.Errorf("%w")` 包装，演示 `errors.Is/As`。
3) 写一个 `safeCall(fn func())`，在当前调用中用 `defer recover` 捕获 panic，打印堆栈。（如果尚未学习 goroutine，可先不启动 goroutine）

## error 处理的常见诟病
- 样板代码多：每步都要 `if err != nil`；使用早返回、包装上下文可以保持简洁。
- 忘记处理错误：忽略返回值或用 `_` 丢弃错误导致静默失败；静态分析工具可检查未处理的 `error`。
- 双重日志/重复包装：多层都记录或包装，导致噪音；应在边界处记录一次，内部只传递上下文。
- 滥用 `panic`：把业务错误当异常；会导致进程崩溃或难以测试。

## error 处理的最佳实践
- 早返回：检测到错误立即返回，保持主流程缩进浅。
- 包装上下文：`fmt.Errorf("fetch user %d: %w", id, err)` 带上关键数据。
- 合理定义哨兵/类型：对可判定的场景使用哨兵或自定义类型，搭配 `errors.Is/As`。
- 明确边界：在 handler/worker 等边界统一记录日志，内部只返回错误。
- 不吞错误：不要用 `_ = fn()` 或空 `default` 吞掉错误，必要时至少记录或显式注释。
- 少用 `panic`：除非不可恢复的编程错误；goroutine 入口使用 `recover` 做兜底。
