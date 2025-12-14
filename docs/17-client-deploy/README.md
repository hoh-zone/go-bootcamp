# 第17章：客户端接入与发布

## 学习目标
- 编写最小客户端验证 SSE 流
- 准备配置与环境变量
- 完成打包与部署检查清单

## 章节提纲
- go build 发布、基础 Dockerfile 思路
- 容器化：构建、传递环境变量、端口映射
- 基础运维检查：日志、健康检查、资源限制

## Docker 与 code/17 示例
- 示例 Dockerfile（见 `code/17/Dockerfile`）：多阶段构建，使用 `golang:1.22` 编译二进制，最终镜像基于 `distroless` 以减小体积、提升安全；暴露 `8082` 端口。
- 构建命令：
```bash
cd code/17
docker build -t chatserver:dev .
```
- 运行命令（需提供大模型访问配置）：
```bash
docker run --rm -p 8082:8082 \
  -e ARK_API_KEY=$ARK_API_KEY \
  -e ARK_MODEL_ID=${ARK_MODEL_ID:-deepseek-v3-250324} \
  chatserver:dev
```
- 环境变量：`ARK_API_KEY` 必填；`ARK_MODEL_ID` 可覆盖默认模型；如使用 JWT 密钥自定义，可在镜像内通过 `JWT_SECRET`（可扩展修改代码）传入。
- 客户端验证：在宿主机打开 `http://localhost:8082`。

## 容器化注意点
- 健康检查：可以在容器编排中配置探针调用 `/healthz`。
- 日志：distroless 输出到 stdout/stderr，交由容器平台收集。
- 资源：为生产环境设置内存/CPU 限制，避免 OOM；为流式接口设置合理超时。
- 安全：尽量使用非 root 镜像（示例使用 `nonroot` 用户）；不要把密钥写入镜像，改用环境变量或密钥管理服务。
