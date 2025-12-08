# HTTP 服务基础示例

对应第 13 章作业，包含：
- `/hello`：查询参数 `name`，默认 gopher。
- `/echo`：POST JSON 回显，支持取消与错误处理。
- `/healthz`：健康检查。
- 中间件：日志 + panic 恢复，`Chain` 组合。
- `NewServer`：封装超时配置。
- `cmd/httpserver/main.go`：提供运行入口，监听 `:8080`。

## 运行
```bash
cd code/13
go test ./...
# 运行示例服务
go run ./cmd/httpserver
```
