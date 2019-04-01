# nautilus 

nautilus是基于golang语言的微服务开发框架。源于国内流行的beego Web开发框架，对go-kit中的微服务工具集合进行了深度集成；另外还集成了基于viper的配置管理，支持将外部consul或者etcd作为配置中心统一配置已下发。

功能包括：

- 服务注册与发现
- 客户端负载均衡
- 熔断
- 统一配置管理
- 错误输出框架

接下来会将继续完善，并添加A+P模式的leader选举，以及A+A模式的共享锁等服务...


1. debug code, product latest container image:
``` 
make all 
```

2. release code, product release container image:
```
make release
```
