- config 配置中心
```go
好的,业务代码调用Config接口获取配置的时候,使用Kratos的config组件后,主要变化是:

不需要关心配置来源,可以统一通过config接口获取,实现解耦。
例如从etcd获取配置:

// 之前的代码

client, err := clientv3.New(...)
configVal, err := client.Get(...)

// 使用kratos config后

c := config.New(config.WithSource(
    etcd.NewConfig(client),
))

val := c.Value("key")
可以同时加载多个配置源,实现配置聚合。
例如同时从Kubernetes和etcd获取配置:

c := config.New(
    config.WithSource(
        kubernetes.NewConfig(k8sClient), 
        etcd.NewConfig(etcdClient),
    ),
)
kratos会把多个源的配置合并后返回。

支持动态配置更新
通过config的Watch方法可以监听变更:

c.Watch("key", func(key string, value config.Value) {
    // 动态更新配置 
})
这样就可以自动响应配置源的变更。

总结一下,kratos config实现了配置抽象和聚合,可以实现从不同源无缝集成配置,不需要变更业务代码。
```

- kratos-config 其他配置项介绍
```go
根据查看kratos源代码,kratos主要封装和适配了以下几种配置来源:

file: 从文件如JSON、YAML等加载配置
env: 从环境变量加载配置
etcd: 从etcd加载配置
nacos: 从Nacos加载配置
consul: 从Consul加载配置
zookeeper: 从Zookeeper加载配置
kubernetes: 从Kubernetes的ConfigMap和Secret资源加载配置
这些配置加载适配都在config目录下的子包中实现,并且都是实现了config.Source接口。

在使用时,可以通过config.WithSource方法添加这些Source,实现从不同数据源加载配置,例如:

import "github.com/go-kratos/kratos/v2/config"

// 创建config
c := config.New(
  config.WithSource(
    file.NewSource("config.yaml")
  ),
  config.WithSource(
    env.NewSource() 
  )
)
kratos会从文件、环境变量等源加载配置,并做合并后返回。

config组件还支持:

监听配置变更事件
原子更新配置
多种格式如JSON、TOML、YAML、Properties
所以kratos封装了常见的配置源,可以通过config统一访问不同来源的配置。
```