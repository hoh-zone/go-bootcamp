# 函数练习示例

实现了本章小练习：
- `Average(nums ...float64) (float64, error)`：输入为空时报错。
- `Filter[T](xs []T, pred func(T) bool) []T`：按谓词过滤切片。
- `NewCounter(start int) func() int`：闭包计数器，每次调用返回递增值。

## 运行
```bash
go run .
```
