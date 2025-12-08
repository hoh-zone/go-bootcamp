# 第2章：环境准备与工具链

## 学习目标
- 安装配置 Go 与常用工具链
- 理解 GOPATH 与 Go Modules 区别
- 完成 go env 检查并创建工作目录
- 运行第一个 go run 示例

## 章节提纲
- 安装 Go、设置环境变量(go env)
- 选择编辑器与插件
- 使用 go mod init 创建模块
- 编译/运行/格式化(go build/run/fmt)

## 安装 Go

### 什么是环境变量，为什么需要它
- 环境变量是操作系统为进程提供的键值对配置，影响程序查找路径、网络代理、语言等行为[1]。
- Go 常见变量：`PATH` 决定终端能否直接运行 `go`；`GOPATH` 指定本地依赖缓存与工作区；`GOROOT` 指向 Go 安装目录（一般自动设置）；`HTTP_PROXY` 等代理变量影响拉取依赖[2][3]。
- 正确设置的好处：命令可直接执行、依赖能顺利下载、模块模式行为可预测；错误设置会导致 `go: command not found`、无法获取依赖或编译路径异常[2][3]。
- 查看当前设置：`go env`（Go 相关）或 `echo $PATH` / `echo %PATH%`（系统 PATH）[3]。


### macOS
- 推荐使用 Homebrew：`brew install go`；升级用 `brew upgrade go`。
- 或从官网下载 pkg 安装包（https://go.dev/dl/），双击安装后默认写入 `/usr/local/go`。
- 安装完成后执行 `go version` 与 `go env GOPATH` 确认；如未自动加入 PATH，可在 `~/.zprofile` 写入 `export PATH=$PATH:/usr/local/go/bin`。

### Windows
- 直接下载官方 MSI 安装包（https://go.dev/dl/），一路 Next 即可，默认写入 `C:\Program Files\Go` 并自动配置 PATH。
- 或使用 winget：`winget install -e --id GoLang.Go`；升级用 `winget upgrade GoLang.Go`。
- 安装完成后打开 PowerShell 运行 `go version`、`go env GOPATH` 检查；如 PATH 异常，手动将 `C:\Program Files\Go\bin` 添加到系统环境变量。

## 环境检查与目录准备
- 运行 `go env` 确认 `GO111MODULE=on`（默认）与 GOPATH 位置。
- 新建练习目录（例如 `~/code/go-class` 或 `C:\code\go-class`），在其中运行 `go mod init example.com/go-class` 以开启 module 模式。
- 验证编译：创建最小示例 `main.go`，使用 `go run .` 运行；格式化可用 `go fmt ./...`。

### 国内加速配置（可选）
- 配置代理：`go env -w GOPROXY=https://goproxy.cn,direct`。`direct` 作为兜底从原始源拉取。
- 校验和：保持默认的 `GOSUMDB=sum.golang.org` 即可；如需跳过特定私有模块，可设置 `GONOSUMDB=example.com/*`。
- 查看生效：`go env GOPROXY GOSUMDB GONOSUMDB`。

## Git 安装与基础使用
- 安装
  - macOS：`brew install git`，或在 Xcode Command Line Tools 中自带。
  - Windows：下载 Git for Windows 安装包（https://git-scm.com/download/win），按默认选项安装并勾选将 Git 加入 PATH。
- 配置：首次使用设置身份 `git config --global user.name "Your Name"` 与 `git config --global user.email "you@example.com"`；可选设置默认分支 `git config --global init.defaultBranch main`。
- 常用命令：`git clone <repo>` 获取仓库；`git status` 查看变更；`git add .` 暂存；`git commit -m "msg"` 提交；`git pull` 同步远端；`git push` 推送。
- GUI 工具：Windows 推荐 GitExtensions（https://gitextensions.github.io/）；macOS 推荐 SourceTree（https://www.sourcetreeapp.com/），便于可视化管理分支、提交与冲突。
- 建议：在本仓库根目录运行 Git 命令，确保文档与代码一起纳入版本控制；提交前可运行 `go fmt ./...` 与 `go test ./...` 保持状态可用。


## 编辑器与常用工具
- VS Code + Go 扩展：会自动安装 `gopls`、`goimports` 等；也可手动执行 `go install golang.org/x/tools/gopls@latest`。
- 终端常备：`go test ./...` 检查单元测试，`go list -m -u all` 查看依赖更新。
- GoLand（JetBrains）：开箱即用的 Go IDE，内置调试、重构与代码分析。学生/教师可通过 JetBrains 学生计划免费申请教育许可（https://www.jetbrains.com/community/education/），每年续期一次。个人推荐先用 GoLand，上手成本低、Go 体验最完整。
- AI 助手
  - GitHub Copilot：VS Code、JetBrains 插件市场直接安装，登录 GitHub 账号，按需开启 Chat/建议。学生可申请 GitHub Student Pack 获得免费额度。
  - Cursor：下载 Cursor 客户端（https://cursor.sh/），内置 AI 补全与对话，可导入 GitHub Copilot Key 或使用自带额度。
  - Codex（本课程工具链）：按课程说明安装 CLI/编辑器插件，支持在终端内调用 AI 助手并与本地文件交互。建议在本仓库根目录运行，便于读取/修改当前项目。
