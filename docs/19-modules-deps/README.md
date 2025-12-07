# 第19章：模块与依赖管理（进阶）

> 第3章已覆盖基础 `go mod` 命令与语义，这里聚焦实践策略、版本管理、私有依赖与多模块协作。

## 学习目标
- 掌握多模块仓库的协作方式（go.work、replace）
- 制定依赖升级/锁定策略与语义化版本发布
- 配置代理、私有模块与校验策略
- 了解可重复构建要点

## 章节提纲
- 多模块协作：go.work vs replace
- 依赖升级/回退与版本策略
- 语义化版本与导入路径规则
- 代理/私有模块与校验策略
- 可重复构建的要点

## 多模块协作：go.work 与 replace
- `go.work` 适合多模块协同开发：在仓库根执行 `go work init ./modA ./modB`，解决跨模块本地引用，无需改 `go.mod`。
- `replace` 仍可用于单模块指向本地路径或特定版本，但发布前应移除避免污染下游。
- 选择：长期多模块协作用 `go.work`；临时调试可用 `replace`。

示例：
```bash
go work init ./service ./lib
go work use ./service ./lib
go run ./service        # 自动使用本地 lib
```

## 依赖升级/回退策略
- 升级：`go get foo@v1.2.3`；回退：`go get foo@v1.1.0`；移除：`go get foo@none`。
- 大批量升级可用 `go get -u ./...`，但需谨慎检查破坏性变更。
- 升级后跑测试/linters，记录变更原因与潜在影响。
- `go mod graph` + `go mod why -m foo` 理解依赖来源，避免引入多余库。

## 语义化版本与导入路径
- 遵守 SemVer：破坏性变更发布 v2+。
- Go 的导入路径必须携带主版本号（v2+），如 `example.com/foo/v2`，避免多版本冲突。
- 打标签：`git tag v1.2.3` 后才可被 `go` 工具解析获取。

## 代理、私有模块与校验
- 代理：`GOPROXY` 可多值，如 `https://goproxy.cn,direct`；内网可搭建私有代理。
- 私有模块：设置 `GOPRIVATE=example.com/*`，跳过公共代理与校验数据库；如有需要设置 `GONOSUMDB`。
- 校验：`go mod verify` 确保本地缓存与 `go.sum` 一致；必要时 `go clean -modcache` 清理。

## 可重复构建的要点
- 依赖锁定：`go.mod` + `go.sum` 即是锁文件，保持提交。
- 配置确定性依赖源：设置合适的 `GOPROXY`/`GOSUMDB`；私有模块用 `GOPRIVATE`。
- 保持一致的 Go 版本与工具：团队内对齐 `go` 版本，避免语义差异。
- 需要离线/隔离时预下载：`go mod download` 预热依赖。

## 调试与常用命令（进阶）
- `go env GOPROXY GOPRIVATE GOSUMDB` 查看关键配置。
- `go version -m <binary>` 检查已构建二进制中的依赖版本。
- `go mod download` 预先下载依赖，便于离线或 CI 预热。

## Makefile 辅助工作区
- 通过 Makefile 把常用的 `go run/test/fmt/tidy` 命令固化，减少记忆成本；用 `.PHONY` 标注伪目标。
- 在多模块工作区下可直接以 `go.work` 为入口运行 `go test ./...`，覆盖所有子模块。
- 推荐的目标：`run`（运行 demo）、`test`（全局测试）、`fmt`（格式化）、`tidy`（逐模块整理依赖）、`work`（`go work sync` 同步）、`clean`（清缓存）、`help`（列出命令）。

示例（`code/19/Makefile`）：
```makefile
.PHONY: run test fmt tidy work clean help

TEXT ?= hello go

run: ## 运行 app，TEXT 控制输入文本
	go run ./app -- "$(TEXT)"

test: ## 在工作区内运行全部测试
	go test ./...

fmt: ## 统一格式
	go fmt ./...

tidy: ## 分别整理 app/lib 依赖
	(cd lib && go mod tidy)
	(cd app && go mod tidy)

work: ## 同步 go.work 与子模块信息
	go work sync

clean: ## 清理构建/测试缓存
	go clean ./...

help: ## 显示可用命令
	@grep -E '^[a-zA-Z_-]+:.*##' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "%-10s %s\n", $$1, $$2}'
```

## 实践作业：多模块协同
目标：用 `go work` 管理一个包含两个模块的工作区，体验本地协作与跨模块调用。
- 模块 A（`lib`）：暴露 `Count(text string) map[rune]int`，统计字符频次。
- 模块 B（`app`）：依赖 `lib`，调用 `Count` 统计输入并打印结果。
- 运行 `go run ./app` 验证本地引用无需修改 import 路径。
