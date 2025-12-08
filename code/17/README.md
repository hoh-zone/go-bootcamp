# SSE 对话后端示例（Chapter 17）

基于第 14 章安全基础，提供：
- `/login`：POST `{username:"alice", password:"123"}`，返回 JWT（Bearer）。
- `/chat`：POST `{message:"你好", model:"your-model-id"}`，鉴权后调用 Ark 大模型流式返回，SSE 输出。
- `/healthz`：探活。
- 中间件：JWT Bearer 校验（跳过 login/healthz）、安全头、日志、recover。
- CLI 客户端：`cmd/chatclient` 交互式登录并循环发起对话，打印 SSE 流，Ctrl-C 退出。

## 运行
```bash
cd code/17
export ARK_API_KEY=...       # 必填
export ARK_MODEL_ID=...      # 模型 endpoint ID，可在请求中覆盖；为空则默认 deepseek-v3-250324
go run ./cmd/chatserver
```

## 调用示例
```bash
# 登录获取 JWT
TOKEN=$(curl -s -X POST localhost:8082/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"alice","password":"123"}' | jq -r .token)

# SSE 聊天
curl -N -X POST localhost:8082/chat \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"message":"你好"}'

# 交互式客户端
go run ./cmd/chatclient
```

## Docker 运行
```bash
cd code/17
# 构建镜像
docker build -t chatserver:dev .

# 运行，记得传入 ARK_API_KEY/ARK_MODEL_ID
docker run --rm -p 8082:8082 \
  -e ARK_API_KEY=$ARK_API_KEY \
  -e ARK_MODEL_ID=${ARK_MODEL_ID:-deepseek-v3-250324} \
  chatserver:dev

# 或使用 docker-compose（包含 healthz 探针）
docker compose up --build
```
