# 安全与校验示例

对应第 14 章内容，演示：
- 登录接口返回 JWT Bearer Token（固定账号密码 alice/123），受保护接口校验签名与有效期（healthz/login 例外）
- `/echo` 请求体验证与取消处理，`/hello` 问候，`/healthz` 探活
- 中间件链：日志、recover、防止 nosniff/iframe/CSP，简单 CORS
- `NewServer` 封装超时配置，入口在 `cmd/secure/main.go`

## 运行
```bash
cd code/14
go test ./...
go run ./cmd/secure

# 默认账号密码：alice / 123，登录后获得 JWT，再用 Authorization: Bearer <token> 访问受保护接口
```
