- rpc 服务开发流程
```go
RPC方法和客户端获取服务配置的流程:

服务端

服务端需要实现proto文件中定义的RPC服务接口,例如Metadata服务:

// api/metadata/metadata.proto
service Metadata {
  rpc ListServices(ListServicesRequest) returns (ListServicesReply);
}
自动生成的代码会包含接口和默认实现:

// api/metadata/metadata.pb.go
type MetadataServer interface {
  ListServices(context.Context, *ListServicesRequest) (*ListServicesReply, error)
}

type UnimplementedMetadataServer struct {}
// 默认空实现
func (*UnimplementedMetadataServer) ListServices() {
  // 空方法  
}
我们的服务端需要嵌入默认实现,并实现业务逻辑:

// server.go
type Server struct {
  metadata.UnimplementedMetadataServer 
}

func (s *Server) ListServices(ctx context.Context, req *ListServicesRequest) (*ListServicesReply, error) {
  // 实现业务逻辑
  // ...
  return &ListServicesReply{
    Services: []string{"service1", "service2"}  
  }, nil
}
客户端

客户端需要基于Stub代码调用RPC方法:

// client.go

conn, err := grpc.Dial(target, grpc.WithInsecure())
client := metadata.NewMetadataClient(conn)

reply, err := client.ListServices(ctx, &ListServicesRequest{}) 
获取服务配置可以通过grpc.WithDefaultServiceConfig指定:

conn, err := grpc.Dial(target, 
  grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))
或者从注册中心获取配置:

config := getServiceConfigFromRegistry() // 获取配置
conn, err := grpc.Dial(target, grpc.WithDefaultServiceConfig(config))
通过以上方式,客户端可以通过Stub调用RPC方法,并使用服务配置启用负载均衡和服务发现。
```