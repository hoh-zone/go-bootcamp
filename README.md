# Go 课程仓库：从零基础到 SSE 实时对话后端

本仓库面向对 Go 语言（Golang）感兴趣但没有实战经验的学习者。课程假定学习者具有一些计算机学科基础（如算法/数据结构、操作系统或网络概念）。如果没有这些背景，课程会在适当位置补充入门资源链接。课程共 19 章，目标是在课程结束后能够使用 Go 实现一个基于 SSE（Server-Sent Events）的流式后端服务，并与 deepseek 进行实时对话集成。

## 📌 关于 AI Coding 工具的重要提示

**本课程充分利用 AI 的力量来加速学习与开发：**

- 📚 **课程内容**：本课程大部分文档与讲义由 AI（如 ChatGPT、Claude 等）生成，经过筛选与补充后呈现。这既提高了内容质量，也体现了 AI 在教育中的实际应用。
- 💻 **代码示例**：绝大部分代码示例由我和 AI 共创完成（Codex、Copilot、Cursor、Claude 等）。我们鼓励你在学习过程中也**充分利用 AI coding 工具**，这已成为现代开发的标准实践。
- 🚀 **学习建议**：**强烈推荐**在学习本课程时安装并使用 AI coding 工具：
  - 快速获得代码补全与智能建议
  - 理解代码的同时学习 AI 如何思考
  - 高效完成练习与实战项目
  - 掌握如何与 AI 高效协作——这本身就是未来工程师必备技能

**AI coding 已不是可选项，而是提高生产力的必需工具。在学习中使用它，不仅能加速进度，还能更好地理解现代开发工作流。**

主要结构
- docs：课程文档，按章节组织（详见下方章节列表）
- code：课程示例与练习代码

学习目标（整体）
- 掌握 Go 的基础语法、并发与网络编程
- 理解并实现 SSE 协议并处理并发订阅/广播
- 能够编写测试、处理错误与上下文取消
- 学会将服务打包与部署（包含基础 Docker 思路）

受众与前置条件
- 面向零 Go 基础，但建议具备基础编程/计算机概念
- 无相关背景者：课程中会提供补充阅读链接（算法、HTTP、并发基础等）

## 🎯 课程分阶段学习路径

### **第一阶段：语言基础（第 1-12 章）**
涵盖 Go 的基本语法、并发与工具链。如果你已有编程基础，可以快速略过或选择关键章节重点学习。

### **第二阶段：工程与实践（第 13-19 章）**
通过构建一个**完整的 SSE 流式对话后端**来学习。强调 **Learning by Doing**，即边做边学。

### **💡 推荐学习路径**

#### **路径 A：完整深入学习（推荐新手）**
- 按顺序完成第 1-19 章
- 每章都有对应的代码示例和练习
- 掌握从基础到进阶的完整知识体系

#### **路径 B：快速上手（推荐有编程经验的学习者）**
如果你已有其他语言（如 Python、Java、Rust）的编程基础，可以：
1. **快速浏览** 第 1-12 章（重点关注 Go 特有的特性：指针、接口、goroutine、channel）
2. **直接进入** 第 13 章开始的工程实践
3. 在构建项目的过程中，**按需深入** 基础章节

推荐快速浏览基础的顺序：
- 第 1-2 章：概览与环境（必读）
- 第 3、6、7、8 章：Go 语法基础（快速浏览）
- 第 10-12 章：Context + 并发 + 通道（重点理解）
- **第 13-19 章：工程实战（逐章深入，边做边学）**

课程结构（19 章摘要）
1. 课程概览与目标 — [docs/01-overview/README.md](docs/01-overview/README.md)  
2. 环境准备与工具链 — [docs/02-env-setup/README.md](docs/02-env-setup/README.md)  
3. 基础语法与入口程序 — [docs/03-basics/README.md](docs/03-basics/README.md)  
4. 数组、切片与映射 — [docs/04-collections/README.md](docs/04-collections/README.md)  
5. 控制流与代码风格 — [docs/05-control-flow/README.md](docs/05-control-flow/README.md)  
6. 函数与错误返回 — [docs/06-functions/README.md](docs/06-functions/README.md)  
7. 指针与结构体 — [docs/07-pointers-structs/README.md](docs/07-pointers-structs/README.md)  
8. 方法与接口 — [docs/08-methods-interfaces/README.md](docs/08-methods-interfaces/README.md)  
9. 错误处理与恢复 — [docs/09-error-handling/README.md](docs/09-error-handling/README.md)  
10. Context 与取消 — [docs/10-context/README.md](docs/10-context/README.md)  
11. 并发基础 — [docs/11-concurrency-basics/README.md](docs/11-concurrency-basics/README.md)  
12. 通道与并发模式 — [docs/12-channels-patterns/README.md](docs/12-channels-patterns/README.md)  
13. HTTP 服务基础 — [docs/13-http-server/README.md](docs/13-http-server/README.md)  
14. 安全与校验（Security & Validation） — [docs/14-security-validation/README.md](docs/14-security-validation/README.md)  
15. 对话接口设计 — [docs/15-chat-api-design/README.md](docs/15-chat-api-design/README.md)  
16. 实现 SSE 流式对话后端 — [docs/16-sse-chat-backend/README.md](docs/16-sse-chat-backend/README.md)  
17. 客户端接入与发布 — [docs/17-client-deploy/README.md](docs/17-client-deploy/README.md)  
18. 模块与依赖管理 — [docs/19-modules-deps/README.md](docs/19-modules-deps/README.md)  
19. 测试与基准 — [docs/20-testing/README.md](docs/20-testing/README.md)

快速开始

**如果你是编程新手或零基础学习者：**
1. 阅读第一章与环境章节：  
   - [docs/01-overview/README.md](docs/01-overview/README.md)  
   - [docs/02-env-setup/README.md](docs/02-env-setup/README.md)
2. 按顺序学习第 1-12 章（语言基础，含 Context/并发/通道）
3. 再进入第 13-19 章（工程实践）

**如果你已有其他编程语言基础，想快速上手：**
1. ⚡ 快速浏览第 1-12 章（可按上面"路径 B"的建议选择重点章节）
2. 🎯 **直接从第 13 章开始**：[HTTP 服务基础](docs/13-http-server/README.md) —— Learning by Doing！
3. 按需回顾基础章节中的知识点

**通用步骤：**
- 在 [code](code/) 中打开相应练习目录，按章节运行示例（每章 README 会给出运行说明）。
- 推荐 Go 版本：Go 1.18+。
- 使用 `go test` 运行单元测试，`go build` 打包可执行文件。
- **推荐使用 GitHub Copilot 或其他 AI coding 工具辅助学习与编码**。

对没有计算机学科背景的学习者
- 课程内会链接外部入门资料（算法、网络基础、进程/线程概念等），并在章节中指向补充阅读。
- 建议先阅读：操作系统/网络/数据结构的入门教程，以便更快理解并发与网络编程概念。

输出与项目目标
- 最终产出：一个可运行的 Go 后端，支持 POST 发送消息与 SSE 订阅，测试与简单部署说明。与 deepseek 集成，实现实时对话能力（集成示例与说明位于后期章节）。

仓库位置
- 文档：[docs](docs/)  
- 代码：[code](code/)

贡献与反馈
- 欢迎提交 issue 或 PR，课程会持续更新示例与修正说明。
