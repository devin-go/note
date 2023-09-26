# 分布式系统幂等性
token 机制
```
用户访问界面时，调用服务端获取token
服务端将token保存到redis里
用户提交表单时，带上token
服务端处理表单时，删除token(如果redis里有这个token，删除时会返回1，否则会返回0)
根据删除结果返回1，则继续处理
否则重复提交
```
status或version机制
```
原理都是 update 时加个where条件，如：
update table set aa='aa',status=2 where id=1 and status=1
```
唯一键机制
```
利用数据库唯一键
```
加锁机制
```
即获取锁后，再往下处理
```