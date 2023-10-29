# 什么是bigkey
```
string 过大， 10 kb
集合元素过多，10000个
```
# 危害
```
阻塞工作线程
    因为 Redis 单线程特性，如果操作某个 Bigkey 耗时比较久，则后面的请求会被阻塞。
内存分配不均
    redis 集群环境
过期可能阻塞
    4.0以下，不支持异步删除
```
# 解决思路
```c
拆分成多个key
hash bigkey案例
    每月订单，key:=order_202310, field=orderid, value=aa
    拆分key思路
        orderidCode:=hash(orderid)%1000
        newKey:=order_202310_+orderidCode
        hget(newKey, orderid)
        hset(newKey, orderid, value)
```
# 为什么hash bigkey 会变慢？
```
(1)哈希冲突，redis使用链式解决
(2)虽然使用 增量哈希扩容，但相对比原来还是慢
```