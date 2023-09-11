Endpoint 功能划分

```go
internal/endpoint 包中的 endpoint.go 文件提供了一些处理服务端点 URL 的工具函数:

- NewEndpoint 函数可以根据 scheme 和 host 创建一个新的 url.URL 对象来表示一个端点。

- ParseEndpoint 函数可以从一个端点 URL 字符串列表中解析出使用指定 scheme 的端点的 host 部分。

- Scheme 函数可以根据一个 scheme 和是否启用安全的标志来返回完整的 scheme,如果启用安全则会在 scheme 后面加上 "s"。

主要功能是创建和解析服务端点的 URL,以及处理端点的 scheme。这可以用来在 Kratos 中统一处理服务注册和发现时的端点信息。

例如可以用 NewEndpoint 和 ParseEndpoint 在不同组件之间转换端点数据,Scheme 则可以根据配置来生成实际的 HTTP 或者 HTTPS 端点 URL scheme。
```


