# 质量前移
## 目标
```
提升代码覆盖率
封装一些公共的方法，方便组内同学使用，提升工作效率
```
## mock的范围
```
mock http 服务
mock grpc 服务
mock db
mock redis
mock sdk
```
mock http
```
效果：
    woapi.ReadV3("").WithJson(obj).Assert(obj)
    woapi.ReadV3("").WithJson(obj).AssertStatus(status)
    woapi.ReadV3("").WithJson(obj).AssertWithComparer(comparer)
分三层封装：
    (1)httptest+gin启动本地服务
    (2)封装简化request
        req.WithJson()
        req.WithHead()
        req.SetExpect()
        比较处理结果
    (3)根据业务进一步封装api
        woapi.ReadV3("/session/",md...)
        woapi.GetV3("",md...)
```
mock grpc
```
效果：
    kmgrpc.FileSys.Set(method, values...)
    kmgrpc.FileSys.ClientStream(method, func(msg))
    kmgrpc.FileSys.ServerStream(method, func(msg))
    kmgrpc.FileSys.BidirectionalStream(method, func(msg),func(msg))
设计：
    可以设置返回结果，流可以设置“勾子”函数
实现：
    mock 掉grpc conn（grpc.Dail()方法）
    mock 掉conn.Invoke()
    mock 掉conn.NewStream(),返回FakeStream(实现grpc.Stream接口)
```
mock redis、db
```
直接用gomonkey.ApplyFuncReturn()
```

## 职责
```
mock http 服务
mock grpc 服务
封装公共中间件
封装http request
封装比较结果（用cmp）
```
mock http
```
httptest+gin 本地服务
    WithJson(obj)
    WithForm()
    WithFile(name,file)
woapi
    路由
    中间件(read、write)
```
mock grpc
```
效果：
    kmgrpc.FileSys.Set(method, value)
原理：
(1)mock 掉grpc.Dial返回的连接对象
(2)mock 掉连接的Invke等方法
(3)mock 掉连接的NewStream方法
```
[gomonkey原理](https://juejin.cn/post/7133520098123317256)
