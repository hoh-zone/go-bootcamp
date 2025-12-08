# 测试与基准示例

对应第 20 章，包含：
- `Count` 及其 table-driven 测试、基准、示例。
- `Process` 演示接口驱动设计，测试中用 fake 替代真实实现。
- `//go:generate mockgen ...` 指令展示生成 mock。

## 运行
```bash
go test ./...
go test -bench .

# 可选：安装 mockgen 生成接口 mock
go install go.uber.org/mock/mockgen@latest
```

## go:generate 简介
- 在源文件中添加形如 `//go:generate <command>` 的指令，`go generate ./...` 会执行它们（按文件声明的顺序），常用于生成 mock、代码模板等。
- 生成命令不会在构建时自动执行，需显式运行。
- 示例：
```go
//go:generate mockgen -destination mock_foo_test.go -package foo . FooInterface
type FooInterface interface {
	Do(x int) error
}
```
- 运行 `go generate ./...` 后会生成 `mock_foo_test.go`，可在测试中直接引用。
