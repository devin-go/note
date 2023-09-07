# 质量前移
## 目标
```
提升代码覆盖率
封装一些公共的方法，方便组内同学使用，提升工作效率
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
原理：
(1)mock 掉grpc.Dial返回的连接对象
(2)mock 掉连接的Invke等方法
(3)mock 掉连接的NewStream方法
```
[gomonkey原理](https://juejin.cn/post/7133520098123317256)
