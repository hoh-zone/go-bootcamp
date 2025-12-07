# 第16章：实现 SSE 流式对话后端

> 课前准备：提前在 https://console.volcengine.com/ark/region:ark+cn-beijing/apiKey?apikey=%7B%7D 获取火山方舟 ARK 的 API Key（保存为环境变量 `ARK_API_KEY`），并准备好模型 endpoint ID（如 `deepseek-v3-250324` 或自定义 `ARK_MODEL_ID`）。

## 学习目标
- 实现可并发的 SSE 推送 handler
- 构建消息广播/订阅机制
- 处理客户端断开与资源清理
- 为 SSE 逻辑编写测试与探活

## 章节提纲
- 认证与入口：登录获取 JWT、受保护 SSE/REST 接口（参考 code/17）
- SSE handler：建立流式响应、设置头、心跳/断线处理
- 订阅管理：注册/移除连接、背压与发送缓冲
- 消息来源：调用大模型流式响应并转发到 SSE（Ark API 示例）
- 取消与清理：`context.Done()`、超时、客户端断开时回收资源
- 测试与验证：curl/交互客户端验证、模拟中断、错误路径覆盖

## 交互客户端（chatclient）说明
- 位置：`code/17/cmd/chatclient`。
- 功能：提示输入用户名/密码，调用 `/login` 获取 JWT，随后循环读取用户输入并向 `/chat` 发起请求，实时打印 SSE 返回的 `data:` 增量；Ctrl+C 会通过 context 取消当前请求并退出。
- 使用方式：先运行服务 `go run ./cmd/chatserver` 或 Docker/Compose 启动；另开终端执行 `go run ./cmd/chatclient` 按提示操作。
