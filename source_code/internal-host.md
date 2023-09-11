host 功能划分
```go
 internal/host/host.go 文件。

通过查看代码,internal/host 包主要提供了以下功能:

1.提取地址中的主机名和端口号:
func ExtractHostPort(addr string) (host string, port uint64, err error)
2.检查IP地址是否有效:
func isValidIP(addr string) bool
3.获取监听器的端口号:
func Port(lis net.Listener) (int, bool) 
4.从地址中提取私有IP和端口:
func Extract(hostPort string, lis net.Listener) (string, error)
它通过解析地址字符串,获取网络接口信息等方式,来实现获取主机名、端口号、IP地址以及从地址中提取私有网络IP的功能。

这些函数被服务发现和注册组件用来处理服务地址信息。

总体来说,internal/host 主要是提供了一些地址处理和解析的工具函数,用于获取地址相关信息。

```