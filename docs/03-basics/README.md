# 第3章：基础语法与入口程序

## 学习目标
- 熟悉 main 包与入口函数
- 掌握包、导入与作用域规则
- 理解变量、常量、短变量声明
- 认识基础类型与打印调试

## 章节提纲
- 变量/常量、类型推断、短变量声明
- 包与模块、可见性规则
- 基本类型与零值、类型转换
- 字符串与 rune/byte 区别
- fmt/print/log 用法与简单调试


## 变量与常量
- 变量声明：`var name type`；可选初始化：`var count int = 3`；批量：`var a, b = 1, "hi"`。
- 零值：整数/浮点为 `0`、布尔 `false`、字符串 `""`、指针/切片/映射/函数/接口为 `nil`，无需手动赋默认值。
- 常量：`const name = value`，值必须在编译期可确定（数字、字符串、布尔、rune）；不可使用 `:=`；常见配合 `iota` 自增生成枚举值。
- 类型推断：有初始值时可省略类型，由编译器推断（但不会跨包推断）。

```go
const (
	StatusOK = 200
	StatusCreated
)

var (
	port int    // 0
	host string // ""
)
```

## 短变量声明（:=）
- 只能在函数内使用；语法 `name := expr`，自动推断类型。
- 至少有一个新变量才允许：`x := 1; x, y := 2, 3`（y 为新变量）；否则编译报错。
- 常见搭配：`v, err := someCall()`。


## 基本类型与类型转换
- 数值：`int`/`uint`（与架构位数相关）、`int8/16/32/64`、`float32/64`、`complex64/128`。
- 布尔：`bool` 只能是 `true/false`，不允许与整数混用。
- 字符：`byte` 是 `uint8` 的别名，`rune` 是 `int32` 的别名，常用于区分按字节或按 Unicode 码点处理。
- 强制转换：Go 不做隐式类型转换，不同整数/浮点类型需显式转换；字符串与数字需配合 `strconv` 或 `fmt`。

| 类型           | 说明                         | 零值   | 常见场景                |
| -------------- | ---------------------------- | ------ | ----------------------- |
| `bool`         | 布尔                         | `false`| 条件判断、开关          |
| `string`       | UTF-8 字节序列               | `\"\"` | 文本、日志              |
| `int`/`uint`   | 与架构位数相关的整数         | `0`    | 计数、长度、循环        |
| `int64`/`uint64`| 固定 64 位整数              | `0`    | 时间戳、ID、金额分      |
| `float64`      | 浮点                         | `0`    | 小数计算、比例          |
| `complex128`   | 复数（两份 float64）         | `0+0i` | 科学计算                |
| `byte`         | `uint8` 别名，单个字节       | `0`    | 原始数据、编码          |
| `rune`         | `int32` 别名，Unicode 码点   | `0`    | 处理字符、遍历字符串    |
| `[]T`          | 切片，引用类型               | `nil`  | 可变长度序列            |
| `map[K]V`      | 映射，引用类型               | `nil`  | key-value 存取          |

```go
var i int32 = 42
f := float64(i)
s := fmt.Sprintf("%d", i) // 数字转字符串
n, err := strconv.Atoi("123")
```

## 字符串、rune 与 byte
- 字符串是只读的 UTF-8 字节序列，`len(str)` 返回字节长度，索引 `str[i]` 得到单个 `byte`。
- `[]byte(str)` 获得底层字节切片；`[]rune(str)` 以 Unicode 码点切分，常用于处理多字节字符长度。
- `for _, r := range str` 以 rune 遍历，避免中文等多字节字符被截断；若按字节遍历，使用传统 `for i := 0; i < len(str); i++ {}`。

```go
s := "Go语言"
fmt.Println(len(s))        // 8 字节
fmt.Println([]byte(s))     // [71 111 ...]
fmt.Println(len([]rune(s))) // 3 个字符
for i, r := range s {
	fmt.Printf("%d -> %c\n", i, r)
}
```

## 入口程序与 main 包
- Go 应用的入口是 `package main` + `func main()`，`go run .` 或 `go build` 会寻找同级目录下的 `main` 包并执行 `main()`。
- 目录即包：同一目录内文件的 `package` 声明必须一致；模块路径(`module example.com/hello`) 只在导入时使用，包名通常为目录名。
- `main` 包只能被编译成可执行文件，不能被其他包导入；库代码放在非 main 包目录下再被导入。
- `init()` 用于做轻量初始化，每个文件可有多个，执行顺序：导入包的 `init()` → 当前包各文件自上而下的 `init()` → 最终执行 `main()`。常见用途：注册驱动、设置日志、校验环境变量，也可以在这里提前解析命令行 `flag`。

```go
// code/02/main.go
package main

import (
	"flag"
	"fmt"
)

var name = flag.String("name", "gopher", "who to greet")

func init() {
	flag.Parse()
	fmt.Println("init runs before main, name flag =", *name)
}

func main() {
	fmt.Println("hello,", *name)
}
```

## 包导入与作用域
- 导入写在 `import` 之后，单行或分组都可以；分组更便于管理第三方依赖和标准库。
- 可选用别名：`import m "math"`；只执行副作用用空白标识符：`import _ "net/http/pprof"`。
- 包可见性：标识符首字母大写即导出（对其他包可见），小写仅包内可用；与文件名、目录名无关。
- 作用域：`package` 级（同包文件共享）、`file` 级 `init`/`var`、代码块级（`if/for/func` 等）。避免变量遮蔽（重名导致外层变量被屏蔽）。
- 

```go
package main

import (
	"fmt"
	mathAlias "math"
)

var pkgLevel = "package scope"

func main() {
	msg := "block scope"
	fmt.Println(pkgLevel, msg, mathAlias.Pi)
}
```

只执行副作用的导入示例：空白标识符 `_` 会运行被导入包的 `init()`，但不会将其中标识符带入当前作用域。

```go
package main

import (
	_ "net/http/pprof" // 注册 pprof handler 到默认 mux
	"log"
	"net/http"
)

func main() {
	log.Println("pprof ready on http://localhost:6060/debug/pprof")
	_ = http.ListenAndServe("localhost:6060", nil)
}
```

## go.mod 是什么
- `go.mod` 记录当前模块的**模块路径**、最低 Go 版本以及依赖列表。根目录一个模块，多模块仓库则子目录各有 `go.mod`。
- 本仓库示例：`code/02/go.mod` 内容：

```go
module example.com/go-class
go 1.22.0
```

- 关键字段：
  - `module <path>`：导入时的前缀，通常对应你的代码托管地址；本地使用时可以是自定义路径。
  - `go <version>`：声明最低兼容的 Go 版本，决定语法/标准库行为。
  - `require <path> <version>`：列出直接依赖；`indirect` 标记表示通过其他依赖间接引入。
  - `replace <old> => <local-or-new>`：开发时指向本地目录或替换版本，常用来调试未发布的模块。
- 常用命令：`go mod tidy` 清理/补全依赖；`go list -m all` 查看依赖树；`go env GOPATH` 看到缓存位置（下载的依赖存放在 GOPATH/pkg/mod）。
- 与其他语言包管理的区别：
  - 无中央 registry 强绑定：Go 通过模块路径（通常是仓库地址）直接拉取源码，不依赖单一私服，类似 `git clone`；可配置私有域名代理。
  - 版本解析简单：语义化版本 + 最小版本选择（MVS），不会出现 NPM 那样的“不同子树装不同版本”或 Python 的“最后安装覆盖”。
  - go.sum 锁定：`go.sum` 记录下载模块版本的校验和，保证可重复构建；无需单独的 lock 文件格式。
  - GOPATH 缓存：下载的依赖放在 GOPATH/pkg/mod，多个项目复用；删除项目不会删除缓存。

## fmt/print 与调试
- `fmt.Print/Printf/Println` 输出到标准输出；常用占位：`%v` 值、`%+v` 展示字段名、`%#v` Go 语法表示、`%T` 类型、`%q` 打印带引号字符串/rune。
- `fmt.Sprintf` 返回字符串便于日志或拼接；`fmt.Errorf` 搭配 `%w` 包装错误。

## Go设计的简洁性
- 包管理为何这样设计：沿用“路径即标识”的思路，直接从远端仓库拉源码，避开单一中心；搭配 MVS + `go.sum` 保证构建可重复且解析规则简单，团队容易理解。
- 大小写导出规则：首字母大写即公开，省去 public/private/protected 关键字，降低语言表面积；调用方读到名称即可判断可见性，统一风格减少认知切换。
- 包名/标识符约定：包名短小、与目录一致、小写无下划线，标识符用驼峰，减少样板与噪音；强调“写起来少、读起来顺”，符合 Go “少即是多”的设计理念。

## 小练习
1) 编写 `main` 包打印运行的当前时间与命令行参数（使用 `flag.Parse()`）。  
2) 读取字符串变量并分别输出其字节长度和 rune 数；用 `range` 打印每个字符。  
