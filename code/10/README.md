# Context 取消与超时示例

对应第 10 章内容，展示：
- `ProcessAll`：消费 job 通道，尊重 `ctx.Done()` 退出，避免泄漏。
- `DoWithTimeout`：用派生的超时 context 包裹操作。
- `WithRequestID`/`RequestID`：用 `WithValue` 传递轻量元数据。

## 运行
```bash
cd code/12
go test ./...
```
