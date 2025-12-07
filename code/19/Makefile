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
