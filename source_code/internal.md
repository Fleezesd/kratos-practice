// internal 简要介绍

// internal 包提供了一些内部使用的类型和功能,不对外暴露
目录结构
```go
.
├── README.md
├── context
│   ├── context.go
│   └── context_test.go
├── endpoint
│   ├── endpoint.go
│   └── endpoint_test.go
├── group
│   ├── example_test.go
│   ├── group.go
│   └── group_test.go
├── host
│   ├── host.go
│   └── host_test.go
├── httputil
│   ├── http.go
│   └── http_test.go
├── matcher
│   ├── middleware.go
│   └── middleware_test.go
└── testdata
    ├── binding
    │   ├── generate.go
    │   ├── test.pb.go
    │   └── test.proto
    ├── complex
    │   ├── complex.pb.go
    │   └── complex.proto
    ├── encoding
    │   ├── test.pb.go
    │   └── test.proto
    └── helloworld
        ├── generate.go
        ├── helloworld.pb.go
        ├── helloworld.proto
        ├── helloworld_grpc.pb.go
        └── helloworld_http.pb.go
```

