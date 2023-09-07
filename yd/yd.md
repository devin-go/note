# 工作内容总结
```
	统一服务启动
		目的：支持服务启动设置一些公共参数和方法
		http:
			公共功能：注册到网关、崩溃恢复、超时设置、限流
			要合并：gin、beego 框架
			支持：组合多个服务
		grpc:
			崩溃恢复、超时设置、限流、消息体设置
		还要支持 tcp(封装custom)
		
	链路追踪,日志级别的
		gin、beego、websocket、tcp:收到消息，把traceid设置到ctx里
		跨服务传递：grpc拦截器，发送时从ctx里提取设置到metadata，接收时从metadata设置到ctx
		消息异步传递:发送消息时，从ctx里提取设置到设置到消息体里，接收时从消息里提取设置到ctx
		初始化、后台服务类型的traceid
		异步的ctx处理
		
	http-->grpc:
		效果：根据uri透传到grpc_method
		支持：某些接口的入参和出参特殊定制
		原理：
			(1)根据grpc_method反射拿到入参类型
			(2)然后invoker过去
		
		
	服务之间传送文件：
			用户上传文件，http接收到文件--->grpc 服务(http_reader 写到 grpc_writer)
			用户下载文件，grpc-->http (grpc_reader写到http_witer)
			封装分两步：
				(1)grpc stream 封装成io.reader或io.writer 
				(2)再进一步封装upload()或download()
			细节部分：需支持不同的“grpc file message”对象
	中间件：
		公共路由
		解密中间件、加密中间件、验证身份中间件
		json、text、ws等中间件

	错误模型：

		(1)服务内部传输
		(2)转换http状态码
		(3)业务错误(metadata)
	
	转换器：
	
	合并服务：
		多个实例变为一个
		配置合并
	
	难点：
		gin query bind：默认是根据from tag来绑定的，如何支持json tag？
		protomessage 对象不是根据protoName tag来序列化、int64转为数字不是字符串?
项目结构：
	api、servcie、client、config、internal、dao

token机制:
	(1)登陆成功时，db和redis存一份
	(2)存db的好处，服务重启时，用户不用重新输入账号密码
	(3)token 有效时间30天，居然没有续约机制

第三方accessToken机制：
	客户用aesKey加密，服务端用aeskey解密，如果成功，则认为是合法用户，然后生成accessToken。
	大概过程如下：
	用户可以添加，企业app
	app里有appid和aeskey
	当调某个业务接口时，必须调接口获取accessToken
	然后访问业务接口时，带上accessToken

待需要补充知识点：
```
# 扩展知识点
## 工作台
```
就是一个一个应用的集合
应用可以分为：系统自带的+用户自定义的
```
[leveldb](../db/leveldb.md)
[nsq](../mq/nsq.md)