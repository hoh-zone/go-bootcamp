# 第7章：指针与结构体

## 学习目标
- 理解指针与内存地址的意义
- 使用 struct 组织数据与嵌套
- 熟悉标签与序列化需求
- 掌握零值友好的设计

## 章节提纲
- 指针获取与解引用，new vs 取址符
- struct 字面量、字段命名/导出规则
- 嵌套与匿名字段（轻量继承）
- 结构体标签与 json 序列化
- 零值可用设计与构造函数习惯

## 指针基础
- 取地址 `&x`，解引用 `*p`；指针存放的是变量地址。
- `new(T)` 返回 `*T`，字段为零值；`&T{}` 也返回指针，更常见且可内联初始化。
- 传指针以便函数修改调用方数据，或避免拷贝大对象；但不要过度使用。

示例：
```go
func increment(p *int) { *p++ }

func demo() {
	x := 1
	increment(&x)        // x 变 2
	y := new(int)        // *int，值为 0
	fmt.Println(*y)      // 0
	z := &Point{X: 3, Y: 4}
	fmt.Println(z.X, z.Y)
}
```

## 结构体与字段导出
- 字段名大写表示导出，可被其他包访问；小写仅包内可见。
- 结构体字面量：按字段名初始化更清晰：`User{Name: "Alice", Age: 20}`
- 零值可用：让未初始化的结构体也处于合理状态，减少 panic。

示例：
```go
type User struct {
	Name string
	Age  int
	Tags []string
}

func (u *User) AddTag(tag string) {
	u.Tags = append(u.Tags, tag) // 零值的切片可直接 append
}

// 对比值接收者与指针接收者
func (u User) AddTag1(tag string) {
	u.Tags = append(u.Tags, tag) // 改的是副本，调用方看不到新增
}

func (u *User) ChangeName(name string) {
	u.Name = name // 改指针，调用方被修改
}

func (u User) ChangeName1(name string) {
	u.Name = name // 改副本，调用方不变
}

// 值接收者但修改底层切片元素：若 Tags 共享同一底层数组，会影响调用方
func (u User) MutateTag0(newTag string) {
	if len(u.Tags) == 0 {
		return
	}
	u.Tags[0] = newTag // 修改切片底层数据
}
```

## 嵌套与匿名字段
- 结构体可以包含其他结构体，形成组合。
- 匿名字段（嵌入）可把内部字段“提升”到外层，形成轻量继承。
- 对比 Java：Go 没有类继承/父子层级，也没有构造器或重载；组合优先于继承，通过嵌入和接口实现代码复用与解耦。
- Go 的设计目标是简单可组合：嵌入只是一种语法糖，避免深层继承带来的耦合与脆弱基类问题；接口是隐式实现，减少样板代码。

示例：
```go
type Audit struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Profile struct {
	User  // 匿名字段，提升 Name/Age
	Audit // 冲突时需通过具体字段访问
	Bio   string
}

func demoEmbed() {
	p := Profile{
		User: User{Name: "Bob", Age: 30},
		Bio:  "Gopher",
	}
	fmt.Println(p.Name)        // 提升访问
	p.User.Age = 31            // 也可显式访问
	p.UpdatedAt = time.Now()   // 来自 Audit

    p.AddTag("tag1")
    p.AddTag1("tag2")

    fmt.Println(p.Tags)     //["tag1"] 
}
```

## 结构体标签与 JSON 序列化
- 标签用反引号放在字段后：``json:"name,omitempty"``。
- `omitempty` 在零值时省略；`-` 表示忽略。
- 注意导出字段才能被 `encoding/json` 访问。

示例：
```go
type Post struct {
	ID      int       `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content,omitempty"`
	private string    `json:"private,omitempty"` // 不导出，序列化不到
	Created time.Time `json:"created_at"`
}

func marshalDemo() {
	p := Post{ID: 1, Title: "Hello"}
	b, _ := json.Marshal(p)
	fmt.Println(string(b)) // {"id":1,"title":"Hello","created_at":"..."}
}
```

## 零值友好与构造函数
- 让零值可直接使用：切片/映射默认 nil，可用 `append` 或延迟初始化。
- 提供构造函数封装必需字段校验：`func NewUser(name string) (*User, error)`
- 避免过度使用 setters；直接导出字段或提供方法做校验。

示例：
```go
type Config struct {
	Addr string
	Tags map[string]string
}

func NewConfig(addr string) (*Config, error) {
	if addr == "" {
		return nil, errors.New("addr required")
	}
	return &Config{Addr: addr}, nil // Tags 先为 nil，按需初始化
}

func (c *Config) SetTag(k, v string) {
	if c.Tags == nil {
		c.Tags = make(map[string]string)
	}
	c.Tags[k] = v
}
```

## 小练习
1) 定义 `Product`，包含 `ID`、`Name`、价格 `PriceCents`，添加 `json` 标签并实现 `Discount(percent int)` 方法。
2) 定义 `Employee`，嵌入 `User` 和 `Audit`，增加字段 `Title`，演示字段提升与冲突处理。
3) 写一个构造函数 `NewPost(title string, content string) (*Post, error)`，要求标题非空、内容可选，并补全 JSON 标签。 
