httputil 功能划分
```go
Content-Type 主要用于描述HTTP请求和响应中的媒体类型信息,它由类型和子类型组成。  
定义了一些HTTP常量,如基础content type为application

提供了获取content type的函数ContentType和ContentSubtype内容子类型
```