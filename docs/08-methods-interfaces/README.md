# 第8章：方法与接口

## 学习目标
- 为结构体定义方法（值/指针接收者）
- 通过接口抽象行为与解耦
- 理解接口满足是隐式的
- 避免滥用空接口

## 章节提纲
- 方法集与接收者选择
- 接口定义、隐式实现、接口断言
- 接口 vs 结构体：何时抽象
- 空接口与任何类型、类型断言/类型切换
- 小练习：定义 Printer 接口与实现

## 方法与方法集
- 值接收者：调用时会复制接收者，适用于小对象、只读或不需要修改的场景。
- 指针接收者：共享同一实例，允许修改并避免大对象拷贝；通常与需要维护内部状态的方法配合。
- 方法集：类型 `T` 的方法集包含值接收者方法；`*T` 的方法集包含值接收者和指针接收者方法。

示例：
```go
type Counter int

func (c Counter) Value() int   { return int(c) } // 值接收者
func (c *Counter) Inc()        { *c++ }          // 指针接收者 *c = *c +1
func demoMethodSet() {
	var c Counter
	c.Inc()          // 值也可调用指针方法，编译器自动取址
	fmt.Println(c)   // 1
	fmt.Println(c.Value())
}
```

这里的 type Counter int 定义的是一个“新类型”，底层是 int，但与 int 不同的类型。新类型不会继承原类型的方法，也不会自动满足原类型的接口，需要自己实现。你在上面为 Counter 定义的 Value/Inc 方法只属于 Counter/*Counter，普通 int 无法直接调用。

补充区分：

类型定义（如 type Counter int）：新类型，底层相同但方法集独立。
类型别名（如 type MyInt = int）：只是另一个名字，和原类型完全相同，方法、接口实现都会一致。

## 接口与隐式实现
- 接口只描述行为，没有实现；满足接口是隐式的，无需声明 “implements”。
- 小接口优先（接口隔离原则），暴露最小必要方法集合。
- 接口变量为 `(动态类型, 动态值)`；零值为 `nil`。

示例：
```go
type Printer interface {
	Print() string
}

type User struct{ Name string }
func (u User) Print() string { return "user: " + u.Name } // 隐式满足 Printer

func greet(p Printer) {
	fmt.Println(p.Print())
}
```

## 接口作为泛型约束
- 泛型类型参数需要一个“约束”来限定可用的操作，接口正是主要的约束形式：接口里的方法、内建约束 `~`、并集 `|` 都可以限制可接受的类型。
- `any`（`interface{}`）作为约束意味着“不限制”，而像 `comparable`、`fmt.Stringer` 则限制了可用操作（可比较、可调用 `String()`）。

示例：约束要求实现 `Print()` 的类型才能被传入。
```go
type Printer interface {
	Print() string
}

func MapPrint[T Printer](items []T) []string {
	out := make([]string, 0, len(items))
	for _, v := range items {
		out = append(out, v.Print())
	}
	return out
}
```
`MapPrint` 的类型参数 `T` 受 `Printer` 约束，调用者只能传入实现了 `Print()` 方法的类型，编译期即可检查。

## 类型断言与类型切换
- 断言：`u, ok := any.(User)`；非该类型时 `ok` 为 false，直接断言失败会 panic。
- 类型切换：`switch v := any.(type) { case string: ... }` 便于按动态类型分支。

示例：
```go
func handle(v any) {
	switch x := v.(type) {
	case int:
		fmt.Println("int", x)
	case fmt.Stringer:
		fmt.Println("stringer", x.String())
	default:
		fmt.Println("other", x)
	}
}
```

## 何时使用接口
- 先写具体类型/函数，只有在需要跨实现、易替换、可测试时才抽象接口。
- 接口通常由使用方定义（面向调用者），而非实现方。

## Go 接口为何隐式实现
- 设计目标：减少样板和耦合，调用者只关心行为签名，不需要实现方显式声明“implements”。
- 隐式实现让适配更灵活：旧代码也可满足新接口（只要方法匹配），有利于可插拔和测试替身。
- 代价：重命名/修改方法可能在编译期才发现接口不再满足，需配合良好测试。

接口满足断言示例：
```go
var _ Printer = (*User)(nil) // 编译期检查 *User 是否实现 Printer
```
含义：将 `(*User)(nil)` 赋值给 `Printer` 类型的匿名变量 `_`，如果方法不匹配就会编译错误；运行时不会执行。可放在文件顶部作为静态断言/文档。

## 空接口与 any
- `interface{}`/`any` 可表示任意类型，但会丢失静态类型信息。
- 使用场景：容器、通用日志、解码动态数据；处理时需要断言或切换。
- 不要滥用空接口代替合理的类型设计。

## 小练习
1) 定义 `Printer` 接口，声明 `Print() string` 方法；让 `User`、`Product` 实现并打印不同格式。
2) 写一个函数 `LogAll(ps []Printer)`，遍历并打印其输出。
3) 写一个类型 `Box`，内部持有 `any`，提供 `AsString() string`，使用类型断言或类型切换来处理不同具体类型。 
