# 第5章：控制流与代码风格

## 学习目标
- 使用 if/switch/for 处理分支和循环
- 理解 defer 的资源释放顺序
- 练习错误优先返回的写法
- 了解标签与跳出的使用场景

## 章节提纲
- if/else、switch 类型/表达式
- for 三种形式、break/continue
- defer 与栈顺序、常见用途
- 组合条件与早返回习惯
- 示例：小型输入校验流程

## if 与 switch
- 条件无需括号，`if cond {}`，配合 `else if`。可在条件前写短变量声明：`if v, err := f(); err != nil { ... }`.
- `switch` 默认带 `break`，匹配成功后自动退出；`fallthrough` 可刻意落入下一分支（少用）。
- 表达式 switch：`switch x { case 1,2: ... }`; 类型 switch 用于接口断言：
```go
switch v := anyValue.(type) {
case string:
    fmt.Println("string", v)
case int:
    fmt.Println("int", v)
default:
    fmt.Println("other")
}
```

## for、break、continue
- 单条件：`for i < n {}`；经典三段式：`for i := 0; i < n; i++ {}`；无限循环：`for {}`。
- `break` 跳出当前 for；`continue` 进入下一轮迭代。
- 迭代 map 顺序随机（见上一章），range 切片时 `k,v` 为拷贝，修改需用索引。

## defer 语义
- `defer` 调用在函数返回前按后进先出执行，用于关闭文件、解锁等。
- 参数在 `defer` 语句处求值，注意循环中 defer 会积累；必要时在循环内显式关闭或使用匿名函数包裹。
### panic / recover
- `panic` 触发异常终止，栈向上展开，依次执行已注册的 defer；库代码尽量避免对外暴露 panic，推荐返回 error。
- `recover` 仅在 defer 中有效，用于拦截 panic 并转为 error 或记录日志：
```go
func safeRun(fn func()) (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic: %v", r)
        }
    }()
    fn()
    return nil
}
```
- 不要滥用 `recover` 吞掉错误，限于明确的边界（如 goroutine 包装、服务入口）以防隐瞒 bug。
```go
func read(path string) error {
    f, err := os.Open(path)
    if err != nil { return err }
    defer f.Close() // 确保退出时关闭
    // ...
    return nil
}
```

## 错误优先与早返回
- Go 风格是“错误优先、早返回”，减少深层嵌套：
```go
if err := validate(input); err != nil {
    return fmt.Errorf("validate: %w", err)
}
```
- 组合条件时可提取变量、拆分判断，避免长布尔表达式影响可读性。

## 标签与跳出
- 仅在需要跳出多层循环时使用标签：
```go
outer:
for i := 0; i < n; i++ {
    for j := 0; j < m; j++ {
        if bad(i, j) { break outer }
    }
}
```
- 避免用标签做“goto 风格”的流程，保持结构化。

## 示例：输入校验流程
```go
func validateUser(name string, age int) error {
    name = strings.TrimSpace(name)
    if name == "" {
        return errors.New("name required")
    }
    if age < 0 || age > 130 {
        return fmt.Errorf("invalid age: %d", age)
    }
    return nil
}
```
- 通过早返回快速失败，保持主路径清晰；必要时给错误添加上下文。


## 实际开发
- 使用工具保持风格一致：`go fmt`（或 `gofmt`）用于代码格式化；`go vet` 做静态检查；`golang.org/x/lint/golint` 及社区 linters（如 golangci-lint）统一命名/注释/错误处理等规范。建议在 CI 里统一执行，确保团队一致性。
