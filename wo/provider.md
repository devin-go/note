# 云平台接口解耦工作
## 内部服务
```
背景：原来根据build tag 区分云平台，有点乱
目标：梳理代码，抽象拆分成微服务
效果：FileSys,UserSys,NotificationSys
```
FileSys
```
GetMetaData()
GetVersionFile()
UploadFile()
```
UserSys
```
GetInfo()
GetInfoBatch()
等等
```
NotificationSys
```
CollectEvent()//收集事件
Roaming()//漫游
```
## api层
```
梳理出标准api和平台相关的api
独立部署
```
标准与非标准api举例
```c
//标准
/api/office/file/:id 获取文件信息
/api/office/file/:id/open 打开文件
//非标准
/api/office/audit/file/  审核相关
/api/office/forms 表单相关
```
难点
```c
接口根据业务梳理出后，但路径有交叉，且数量比较多，如何升级:
(1)在nginx里写，特定的api就转发，但接口有几十个，有点多
(2)nginx-->apiserver http proxy 如何解决服务发现负载均衡问题？
(3)kscloudSys 设计微服务
//kscloudSys
(1)apiserver与kscloudSys 用grpc stream流通信
(2)apiserver 将http request 序列化到grpc stream里（req.Write(writer)）
(3)kscloudSys反序列为http reqeust对象 (http.ReadRequest(reader))
(4)gin.Router.ServeHTTP(),这样kscloudSys 就可以保持原来的逻辑了
```

