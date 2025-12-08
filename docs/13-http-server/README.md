# 第13章：HTTP 服务基础

> 前置：假设你已掌握第 10-12 章（Context、并发、通道）。如果还不熟悉 HTTP 基本概念（请求/响应、方法、状态码、头部、Body），建议先浏览一篇 HTTP 入门资料，再回到本章实践。

## 学习目标
- 使用 net/http 构建基础服务
- 编写 HandlerFunc 与路由
- 设置响应头、状态码与 cookies
- 实现简单的中间件链

## 章节提纲
- Server 与 Handler 接口、ListenAndServe
- ServeMux/第三方路由器的选择（概念）
- 读取 query/path/body，返回 JSON/文本
- 中间件模式：日志、恢复、跨域
- 健康检查与基础配置结构体

## Server 与 Handler 基础
- `http.ListenAndServe(addr, handler)` 启动服务；`handler` 可为 `http.Handler` 实现或 `http.HandlerFunc`。
- `http.DefaultServeMux` 是默认多路复用器，`http.HandleFunc`/`Handle` 注册路由；自定义 `ServeMux` 可避免全局污染。
- Handler 签名：`func(w http.ResponseWriter, r *http.Request)`；注意 `r.Context()` 携带取消/超时。

示例：
```go
func hello(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "gopher"
	}
	fmt.Fprintf(w, "hello %s", name)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", hello)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
```

## 读取请求与返回响应
- Query 参数：`r.URL.Query().Get("q")`
- Path 参数：标准库 ServeMux 不支持模板匹配，可用前缀匹配或第三方路由器（chi, gorilla/mux）。
- Body：`io.ReadAll(r.Body)` 或 json 解码；记得 `defer r.Body.Close()`。
- 响应：`w.WriteHeader(status)` 设置状态码；`w.Header().Set("Content-Type", "...")` 设置响应头；`json.NewEncoder(w).Encode(data)` 写 JSON。

## 中间件模式
- 典型职责：日志、恢复（recover panic）、CORS、认证、限流。
- 模板：`func(next http.Handler) http.Handler { return http.HandlerFunc(func(w,r){...; next.ServeHTTP(w,r)}) }`
- 组合：从外到内 wrap，或使用切片迭代包装。

示例：日志 + 恢复
```go
func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "internal error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func chain(h http.Handler, m ...func(http.Handler) http.Handler) http.Handler {
	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}
	return h
}
```

## 健康检查与配置
- 健康检查 endpoint：`/healthz` 返回 200/ok，便于探活。
- 配置结构：服务端口、超时、CORS 允许来源等；可用环境变量/标志加载。
- 服务器超时：`http.Server{ReadTimeout, WriteTimeout, IdleTimeout}`，防御慢速攻击。

## 第三方路由器（概念）
- 选择标准：易用性、性能、中间件生态。常见：`github.com/go-chi/chi`, `github.com/gorilla/mux`。
- 标准库足够时优先用标准库；需要 path 参数、分组中间件时可选第三方。

## 小作业
1) 实现 `/hello` 与 `/echo`（POST JSON 回显）的 Handler，并返回合适的 Content-Type/状态码。  
2) 写一个日志中间件和 recover 中间件，组合后挂到 ServeMux 上。  
3) 为 server 添加 `/healthz` 健康检查，配置读/写/idle 超时，并使用 `r.Context()` 支持取消。
