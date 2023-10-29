# 启动服务
```
    http
        设置默认api超时间
        设置bbr限流
        捕获panic
        支持http handler
        向网关注册路
    grpc
        默认超时时间
        消息体大小
        拦截器(traceid)
        捕获panic
```
# 配置
```
支持多种数据元，使用的是etcd
支持热更新，监听某个配置的变化
```
# 错误模型
```go
原来模型
    type Status struct {
        Code int    `json:"code"`
        Msg  string `json:"message,omitempty"`
        At   string `json:"createdAt,omitempty"`
    }
    原来不支持多级错误
    只有业务code，没有对应的http状态
    code不清析
新模型
    支持多级错误，并兼容原来的错误模型
        使用metadata解决
```
# 元数据信息传递
```
    traceid、服务角色
```
# 服务注册与发现
```
服务启动时，可以注册版本信息，用于平滑升级
```
# 路由与负载均衡
```
可以根据版本号选择连接，用于平滑升级
```
# 中间件
```
只使用了
    限流器、链路追踪
未使用
    熔断器、监控等

```