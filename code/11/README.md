# 并发基础示例

对应第 11 章作业，涵盖：
- `ProcessWithPool`：使用 WaitGroup + worker pool 限制并发、保持结果顺序。
- `DoWithTimeout`：`context.WithTimeout` 包裹操作，超时返回错误。
- `SafeCounter`：用互斥锁消除数据竞争。

## 运行
```bash
cd code/11
go test ./...
```
