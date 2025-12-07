# 错误处理练习代码

覆盖本章小练习：
- `ValidationError` 包含 `Field`、`Reason`，实现了 `Error()`.
- `WrapLoad` 调用底层 `load` 并用 `%w` 包装错误，演示 `errors.Is/As`.
- `safeCall` 在当前调用中 `recover` panic 并打印堆栈。

## 运行
```bash
go run .
```
