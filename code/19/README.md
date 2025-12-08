# go.work 多模块示例

本例对应第 19 章作业：
- `lib` 模块提供 `Count(text string) map[rune]int`。
- `app` 模块依赖 `lib`，通过 `go.work` 本地协同。

## 运行
```bash
cd code/19
go run ./app "hello gophers"
```
