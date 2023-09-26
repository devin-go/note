# 与内核交互部分
## 需求(要解决的技术点)
```
1.调用一次内核方法，产生一个结果(需调用 open file，是否成功)
2.调用一次内核方法，可能产生多个结果(如 调用预览文件内容，产生文字，或需要远程下载图片)
3.内核事件，如上报埋点等，即服务端没有调用内核方法，但内核有事件需要上报
4.要支持word，et等内核
```
## 技术方案
### 概括
```
服务与内核属于两个进程通信
且两个进程之间有相互调用的关系，只能用双向流解决
```

### 方案一 用grpc双向流
```
缺点：
需要一个端口(随机端口)
优点：
可以复用grpc 服务，不用自己管理服务的细节
```

### 方案二 socketpair(现在使用的)
socketpair创建一对无名、相互连接的套接字，[原理](https://cloud.tencent.com/developer/article/2169065)
```
fd[0]写的数据，只能由fd[1]读出
而fd[1]写的数据，也只能由fd[0]读出
    go 端使用fd[0]
    core端使用fd[1]
```
# htmlserver
```
(1)index.html 动态注入一些内容：
    wps_env
    file_info
(2)根据不同浏览器做一些差异化：
    如微信小程序进来的，必须登陆
    兼容QQ浏览器不能302问题
(3)预调用内核打开文件
(4)兼容私有化部署：
    如金山文档、钉钉平台，他们之间的对index.html注入的内容有些不一样
```
# 服务自动上下线
```
如何服务内存或CPU使用超过了配置，则下线：
    grpc client 端要支持，随着服务下线，不能再服务新的请求。
实现部分：
    实现grpc resolver
```
# 跨文档复制粘贴
```
copy:
    setRedis(clipboardid, selfRouterInfo)
    selfRouterInfo:
        addr,sessionid,connid
paste:
    routerInfo:=getRedis(clipboardid)
    rpc.Call()
```