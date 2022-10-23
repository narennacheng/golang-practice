该项目参考自[Go进阶训练营](https://u.geekbang.org/subject/go) ，主要用来个人学习go-web框架如何从0开始1的过程。



内容说明，设计是一个迭代的过程：

1. 简单的server、context抽象，自定义的handler处理路由
2. 错误处理与简单路由树实现
3. AOP: 用闭包来实现责任链
  • 为 server 支持一些 AOP逻辑
  • AOP：横向关注点，一般用于解决 Log, tracing，metric，熔断，限流等
  • filter(middleware)：我们希望请求在真正被处理之前能够经过一大堆的 filter
4. 路由树(map路由 -> 路由树) 路由匹配: (全匹配 -> 通配符匹配 -> Todo 复杂匹配)
5. 优雅退出（监听系统信号、 channel select、Hook设计、同步等待与超时机制、context包与线程安全）
6. TODO 静态资源
7. TODO context复用 sync.Pool
8. TODO 接口的注册与发现，Builder。另外可以参考dubbo-go的实现



个别内容待补充..

