# 第6章：函数与错误返回

## 学习目标
- 定义和调用函数（值/指针参数）
- 掌握多返回值与命名返回
- 使用可变参数与闭包
- 理解函数是一等公民的意义

## 章节提纲
- 函数签名、传值语义、指针参数
- 多返回值与错误返回规范
- 可变参数与 ... 展开
- 匿名函数与闭包，函数作为参数
- 纯函数 vs 有副作用的选择

## 函数签名与调用
- Go 一律按值传递；想要修改调用方数据用指针参数：`func set(n *int) { *n = 1 }`
- 函数可以返回多个值，通常把数据与 `error` 放在一起返回：`func load(id int) (User, error)`
- 使用命名返回值时，建议只在简短函数或需要 `defer` 修改返回值时使用，否则保持未命名更清晰。

示例：
```go
// 传值：调用者不会被修改
func addCopy(n int, delta int) {
	n += delta // 只改副本
}

// 传指针：调用者会被修改
func addPtr(n *int, delta int) {
	*n += delta
}

// 命名返回值可配合 defer 微调结果
func loadName(id int) (name string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("load %d: %w", id, err)
		}
	}()
	switch id {
	case 1:
		return "alice", nil
	case 2:
		return "bob", nil
	default:
		return "", errors.New("not found")
	}
}
```

## 错误返回规范
- 返回 `(T, error)` 时，调用方先判错再用值：`v, err := f(); if err != nil { return err }`
- 错误信息应描述失败原因，并携带关键数据：`fmt.Errorf("load %d: %w", id, err)`
- 避免用 `panic` 作为普通错误处理；`panic` 仅用于不可恢复的编程错误。

示例：
```go
func openConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read %q: %w", path, err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse %q: %w", path, err)
	}
	return cfg, nil
}
```

## 可变参数与展开
- 可变参数：`func sum(nums ...int) int`；调用时 `sum(1, 2, 3)`。
- 将切片展开传入：`sum(vals...)`，注意空切片也可展开。
- 可变参数底层就是切片，修改会影响同一底层数组；若不希望被修改，复制一份。

示例：
```go
func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

vals := []int{1, 2, 3}
fmt.Println(sum(vals...)) // 6
```

## 匿名函数与闭包
- 函数是“一等公民”，可赋值给变量、作为参数或返回值。
- 匿名函数可立即调用：`func(msg string) { fmt.Println(msg) }("hi")`
- 闭包捕获外部变量的引用，循环中使用时要注意变量复用问题，例如在 `for` 中传参避免共享同一迭代变量。

示例：
```go
// 立即调用
func() { fmt.Println("once") }()

// 返回闭包
func newCounter(start int) func() int {
	i := start
	return func() int {
		i++
		return i
	}
}

// Go 1.22+ 的 for range 每次迭代都是新变量，无需再复制 i
var fns []func()
for i := range []int{0, 1, 2} {
	fns = append(fns, func() { fmt.Println(i) })
}
fns[0]() // 0

```

## 纯函数 vs 有副作用
- 纯函数：仅依赖输入，无外部状态，重复调用返回相同结果（便于测试与推理）。
- 有副作用的函数：读写 IO、全局变量、时间等。尽量将副作用集中、隔离，方便测试和替换。

示例：
```go
// 纯函数
func add(a, b int) int { return a + b }

// 有副作用：读取时间
func stamp(msg string) string {
	return fmt.Sprintf("%s @ %s", msg, time.Now().Format(time.RFC3339))
}
```

## 函数类型与高阶函数
- 可以为函数定义类型，方便复用签名：`type Pred func(int) bool`
- 高阶函数：参数或返回值是函数。

示例：
```go
type Pred[T any] func(T) bool

func Filter[T any](xs []T, p Pred[T]) []T {
	var out []T
	for _, x := range xs {
		if p(x) {
			out = append(out, x)
		}
	}
	return out
}

evens := Filter([]int{1, 2, 3, 4}, func(n int) bool { return n%2 == 0 })
fmt.Println(evens) // [2 4]

```

## 小练习
写一个函数 `Average(nums ...float64) (float64, error)`，当输入为空时返回错误。
