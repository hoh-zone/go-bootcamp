# 第14章：安全与校验（Security & Validation）

本章目标：讲解在构建 Go 后端时常见的安全注意点与输入校验方法，帮助你在实现 SSE 流式服务和对话接口时减少常见漏洞与错误输入导致的问题。

本章涵盖内容（摘要）

- 输入校验：请求体、URL 参数与 Header 的验证策略、错误返回设计
- 身份认证与授权概念：介绍 JWT 与 API Key 的基本使用与典型实现注意事项，说明在 SSE/长连接场景下的认证续期与断线处理策略。
- 防御常见攻击：CSRF、XSS、注入（SQL/命令注入）
- 服务韧性与滥用防护：请求超时、重试+指数退避、熔断/舱壁隔离、速率限制等手段
- 敏感信息处理：如何避免在日志中记录敏感数据、配置与密钥管理的基本建议（环境变量、Vault 等简介）。

## 输入校验策略
- 入口前置：在 Handler 层对 body/query/header 做白名单验证，拒绝未知字段，返回一致的错误结构。
- Schema 校验：使用 `go-playground/validator` 对 struct tag（如 `validate:"required,email"`）校验，避免手写重复逻辑。
- 错误格式：返回机器友好且用户可读的信息，如 `{ "error": "invalid input", "details": ["email is required"] }`；避免泄露内部错误。
- 长连接/SSE：连接前校验 token/key，连接后可周期校验；连接断开时清理会话状态。

## 鉴权与凭证
- API Key：简单易用，放在 Header；需定期轮换，限制作用域与来源。
- JWT：自包含令牌，注意设置短有效期、签名算法（避免 none）、校验 aud/iss；长连接场景可在握手后存储 claims。
- 会话续期：使用刷新 token 或重连时重新验证；SSE/WS 在后台检测过期时主动关闭连接。
- 最小权限：按路由/动作划分权限，后端进行授权检查，而不仅仅是认证。

### JWT 详解与示例（参考 code/17）
- 组成：Header（算法/类型）、Payload（claims）、Signature（签名）。常用 HS256（对称密钥）或 RS256（非对称）。
- 关键 claims：`sub`（主体）、`exp`（过期时间）、`iat`（签发时间）、`aud`/`iss`（受众/签发方）。避免使用 `alg: none`。
- 生命周期：短有效期 + 服务端轮换 secret/密钥；长连接时在握手后缓存 claims，并在过期时断开或要求重连。
- 校验：验证签名 + 过期时间 + 受众/签发方；拒绝解析失败或过期的 token。
- 传输：放在 `Authorization: Bearer <jwt>`；避免写入日志或返回给无关方。

code/17 片段：
```go
// 登录：固定用户 alice/123，签发 30 分钟有效的 HS256 JWT
func LoginHandler(cfg Config) http.HandlerFunc {
    // 解析请求...
    token := issueJWT(cfg.JWTSecret, req.Username) // JWTSecret=服务端密钥
    json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// 校验：跳过 /login 与 /healthz，其余接口要求 Bearer JWT
func BearerAuthMiddleware(secret string, allowlist []string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // 解析 Authorization: Bearer <token>
            claims, err := parseJWT(secret, rawToken)
            if err != nil || claims.Subject == "" {
                http.Error(w, "unauthorized", http.StatusUnauthorized)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```
- 断开策略：在 SSE 处理时如果 `ctx.Done()` 或 JWT 无效，应及时退出并释放资源。


## 防御常见攻击
- CSRF：对有状态请求使用 SameSite Cookie、CSRF Token 或仅允许 Bearer Token；对跨站请求预检校验 Origin/Referer。
- XSS：输出时使用模板自动转义；返回 JSON 时设置 `Content-Type: application/json`; 对用户输入存储时做存储型 XSS 防护。
- 注入：使用参数化查询，永远不要拼接 SQL；对命令执行使用 `exec.CommandContext` 并限制参数。
- 头部安全：设置 `Content-Security-Policy`、`X-Content-Type-Options: nosniff`、`X-Frame-Options`、`Strict-Transport-Security`（仅 HTTPS）。

## 服务韧性与滥用防护
- 超时：为外部依赖设置 context 超时；服务器端设置 Read/Write/Idle timeout 防慢攻击。
- 重试：对幂等操作使用退避重试；非幂等操作避免自动重试或结合幂等键。
- 熔断/舱壁：隔离下游故障，避免拖垮全局；可用现成库或简单的 semaphore 限制并发。
- 限流：令牌桶或漏桶实现；在网关或服务内拦截异常流量。
- 观测：记录安全相关事件（失败登录、权限拒绝、限流命中），方便审计。

## 敏感信息处理
- 日志脱敏：屏蔽密码、token、身份证号等；必要时仅记录 hash/掩码。
- 配置管理：使用环境变量或秘密管理服务（Vault、Secrets Manager）；避免密钥进仓库。
- 传输安全：优先 HTTPS；内部服务也应考虑 mTLS 或专网隔离。


## 外部参考

- go-playground/validator: https://github.com/go-playground/validator
- OWASP Cheat Sheet（安全实践速查表）: https://cheatsheetseries.owasp.org/
- JWT 简介与实践: https://jwt.io/introduction
- HTTP 与安全头概览: https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Headers
