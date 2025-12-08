# 控制流示例与 lint 练习

## 运行
```bash
go run .
```

## Lint 练习
- 安装：`go install golang.org/x/lint/golint@latest`。若安装后提示 `golint: command not found`，请把 `$(go env GOPATH)/bin`（或自定义 `GOBIN`）加入 `PATH`，如 `export PATH="$(go env GOPATH)/bin:$PATH"`。
- 运行：`./lint.sh` 或 `golint ./...`。
- 本例中 `DoWork` 导出但缺少注释，golint 会报错，练习修复或改为不导出。
