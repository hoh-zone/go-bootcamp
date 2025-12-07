# 第20章：测试与基准

## 学习目标
- 编写 table-driven 单元测试
- 运行 go test 并查看覆盖率
- 了解基准测试与示例测试
- 学习轻量 mock 的技巧
- 认识测试的重要性：防回归、支撑重构、建立设计反馈

## 章节提纲
- testing 包与 *_test.go 规则
- 断言模式与 table-driven 例子
- 测试隔离、临时目录、快/慢测试
- go test -cover/-bench 的用法
- 接口驱动的伪实现与假数据
- 第三方工具：`stretchr/testify` 断言/require，`uber-go/mock` 生成接口 mock

## 为什么测试重要
- 防回归：锁住已有行为，升级/重构更安心。
- 设计反馈：迫使接口更清晰、更易用，过难测的设计通常需要简化。
- 文档作用：示例测试（`ExampleXxx`）可作为可执行文档。
- 速度与信心：小步提交 + 快速测试循环。

## testing 包与 table-driven 测试
- 测试文件以 `_test.go` 结尾，函数签名 `func TestXxx(t *testing.T)`。
- table-driven：用切片列出用例，循环执行，减少重复。

示例：
```go
func TestAdd(t *testing.T) {
	cases := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"positive", 1, 2, 3},
		{"zero", 0, 0, 0},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if got := add(tc.a, tc.b); got != tc.expected {
				t.Fatalf("want %d, got %d", tc.expected, got)
			}
		})
	}
}
```

## 覆盖率与基准/示例
- 覆盖率：`go test -cover ./...`，详细查看 `go test -coverprofile=cover.out` + `go tool cover -html=cover.out`.
- 基准：`func BenchmarkXxx(b *testing.B)`；在循环中调用待测函数，命令 `go test -bench .`。
- 示例：`func ExampleXxx()`，输出用 `// Output:` 注释，既是文档又能被测试执行。

## 测试隔离与临时资源
- 使用 `t.TempDir()` 创建隔离目录；测试完成自动清理。
- 避免依赖全局状态；必要时在 `TestMain` 中集中初始化/清理。

## 断言库与 testify
- 标准库即可完成测试，但断言库可提升可读性。
- `github.com/stretchr/testify/assert` 提供非终止断言；`require` 在失败时立即终止。
- 用法示例：
```go
import "github.com/stretchr/testify/require"

func TestCount(t *testing.T) {
	got := Count("aba")
	require.Equal(t, map[rune]int{'a': 2, 'b': 1}, got)
}
```
- 使用前需 `go get github.com/stretchr/testify@latest`，注意仅作为测试依赖。

## 接口驱动的伪实现与 mock
- 倾向自写假实现（fake/stub）满足接口，保持简单可读。
- 需大量交互断言时可用生成工具，如 `go.uber.org/mock`：
  - 安装：`go install go.uber.org/mock/mockgen@latest`
  - 生成：`mockgen -destination mock_foo_test.go -package foo . FooInterface`
- 在测试中注入 mock 以断言调用、配置返回值。

## AI 辅助开发下的约束
- 测试先行：先写或补齐测试，再让 AI 生成实现，用测试验证结果。
- 明确输入输出：在提示中写清接口契约、边界条件，并声明“不新增依赖”或限定允许的包。
- 自动校验：结合 `go test ./...`、`go vet`、lint 等，让 AI 产物必须通过检查。

## 小练习
1) 为第 10 章的 `Count` 编写 table-driven 测试，覆盖空串、重复字符等。
2) 为一个接口编写手工 fake，实现最小行为并用于测试。
3) 体验 `testify/require`：改写一个测试使用 `require.Equal`，对比标准库写法。
4) 选做：安装 `mockgen` 为某接口生成 mock，在测试中验证调用次数或参数。

