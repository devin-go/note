# 内部数据结构
## dict
```
类似hashmap，但有点不同的时，当内部需要扩展空间的时候(旧数据需要迁移)，不是一次迁移完成，而是把迁移分散到多个get、set等操作上，redis将这个步骤称为“增量式重哈希”，目的是提高响应效率
```
## sds
全称是 Simple Dynamic String（简单动态字符串）  
redis 的定义
```c
struct __attribute__ ((__packed__)) sdshdr8 {
    uint8_t len; /* used */
    uint8_t alloc; /* excluding the header and null terminator */
    unsigned char flags; /* 3 lsb of type, 5 unused bits */
    char buf[];
};
//当然还有 unint16、unint32、unint64类型
```
为什么redis，不直接用c语言的字符串类型(char [] + string.h)呢？  
```
(1)c语言获取字符串的长度，需要每次遍历字符串的，效率低
(2)防止字符串拼接溢出：c语言strcat拼接两个字符串，如果长度不够会溢出，而sds的append函数，会判断len，不够时会扩容
(3)减少申请字符串空间的次数：字符串修改操作的时候，会判断len，如果当前buf够用，就不用申请了
```
## robj
全称是redis object
```c
typedef struct redisObject {
    unsigned type:4;       // 数据类型  integer  string  list  set
    unsigned encoding:4;
    unsigned lru:LRU_BITS; /* LRU time (relative to global lru_clock) or
                            * LFU data (least significant 8 bits frequency
                            * and most significant 16 bits access time). 
                            * redis用24个位来保存LRU和LFU的信息，当使用LRU时保存上次
                            * 读写的时间戳(秒),使用LFU时保存上次时间戳(16位 min级) 保存近似统计数8位 */
    int refcount;          // 引用计数 
    void *ptr;              // 指针指向具体存储的值，类型用type区分
} robj;
```
## ziplist
ziplist是一个特殊的双向链表，用一块连续空间来表示双向链接，如
```xml
<zlbytes><zltail><zllen><entry>...<entry><zlend>
<entry>的结构:<prevrawlen><len><data>
```
充分体现了redis对存储空间的极致要求，存储样例
![Alt text](image.png)
### ziplist与hash
当hash的key-val 数量比较少时，底层就是用ziplist来存储的，如果key-val的数量 或 val 的值超过这两个配置，那么hash的底层就会转会dict来存储
```
hash-max-ziplist-entries 512
hash-max-ziplist-value 64
```
## quicklist
对应的就是redis list数据类型(类似golang的切片)，支持向头部和尾部增加、删除元素，且这些操作都是O(1)复杂度
```c
//O(1)复杂度
lpush、lpop、rpush、rpop
//O(N)复杂度，任意中间位置的存取操作
lindex、linsert
```
### 底层结构
```
双向链表+ziplist，节点是ziplist  
```
为什么要这样设计？
```
(1)如果是只是双向链表，会产生很多的内存碎片
(2)如果只是用类似sds结构来表示，会有扩容时拷贝问题
综上所述，用双向链表可以解决扩容时拷贝问题、用ziplist可以尽量减少内存碎片问题
```

## sorted set
redis的有序集合，例子
```c
//ranking 集合名称 80分数，a 是名字
redis 127.0.0.1:6379> ZADD ranking 80 a
(integer) 1
redis 127.0.0.1:6379> ZADD ranking 90 b
(integer) 1
redis 127.0.0.1:6379> ZADD ranking 85 c
(integer) 1
//还有排序，范围查询等命令 https://www.runoob.com/redis/redis-sorted-sets.html
```
当数据比较少时，sorted是由ziplist实现的，当数据比较多时，则由skiplist+dict实现，具体可以根据配置控制
```c
zset-max-ziplist-entries 128 //当sorted set的元素大于128个时，转为skiplist+dict
zset-max-ziplist-value 64 //当数据项(例子中的名字)大于64个字节时，转为skiplist+dict
```
### skiplist 原理
首先是一个有序的链表，但在有序链表上操作的复杂度都是O(n)，如果能像有序数组那样用二分查找就好了(引出level)
即skiplist是[有序链表+level构成](https://www.bilibili.com/video/BV1tK4y1X7de/?spm_id_from=333.337.search-card.all.click&vd_source=8db508e686bc7c5ece7e6b27eabf33e5)
![Alt text](image-1.png)

构造过程从L1层开始(有多少层是一有个算法的)  
查找过程从L3层开始  
层数算法
```c
randomLevel()
    level := 1
    // random()返回一个[0...1)的随机数
    while random() < p and level < MaxLevel do
        level := level + 1
    return level
/*
redis的参考值
p = 1/4
MaxLevel = 32
*/
```

## HyperLogLog
HyperLogLog是一种算法，典行的场景，用于统计日活、用活等，如统计google主页的日活  
//注意：不能拿到某个用户当天的访问次数
```c
//redis的hyperloglog只支持三个命令
PFADD key [element [element ...]]
PFCOUNT key [key ...]
PFMERGE destkey sourcekey [sourcekey ...]
```
[详细介绍](https://juejin.cn/post/7158403890826706981)

# 面试
[源文](https://www.xiaolincoding.com/redis/base/redis_interview.html)
## redis简介
```
redis是一个内存数据库，读写操作都在内存完成，因此速度非常快，常用于缓存、分布式锁等场景
redis提供string、hash、list、set等常用数据类型
除此外redis还支持事务、lua脚本
还支持将数据持久化到磁盘
支持主从、哨兵等集群模式
```
## redis为什么这么快？
```
(1)redis 的大部分操作都在内存完成
(2)redis采用单线程模型避免了多线程之间的竞争，省去了多线程上下文切换
(3)redis使用io多路复用
```
## redis是单线程吗？
```
redis单线程是指处理用户发的命令这一流程是单线程[接收命令-->解析命令-->处理-->发送结果]，称为redis 主线程 
除了主线程外，还有持久化线程、lazyfree线程(4.0后增加)
```
## redis6.0后为什么引用多线程？
redis的瓶颈并不是cpu，而更多的是内存和网络io,所以单线程并没有什么问题。  
6.0之后引入的是多个线程处理网络io，默认只开启了发送数据多线程，接收处理请求没有开启(真正处理命令还是单线程的)  
### 接收命令开启多线程
```c
//读请求也使用io多线程
io-threads-do-reads yes 
// io-threads N，表示启用 N-1 个 I/O 多线程（主线程也算一个 I/O 线程）
io-threads 4
```

## redis持久化
有三种持久化方式：AOF、RDB、混合持久化(RDB+AOF)
### AOF
记录写的命令到文件里，流程大概是这样：
```c
(1)客户端发送写命令
(2)redis执行写命令
(3)将命令写到aof_buf 缓存区
(4)调用系统io.write将aof_buf缓存的数据写到文件里
(5)第四步其实只是将数据写到系统page cache里，并未真正持久到磁盘
//注意是先执行命令，然后再写到文件和数据库的先写文件策略不一样
```
命令持久化到磁盘有三种策略(控制的是第5步)
```
1.每执行一次写命令，就持久化到磁盘
2.每秒持久化一次
3.no 由系统决定
```
#### AOF重写机制
如果一直把写命令记录到文件里，文件肯定会非常大，所以有一个重写机制，将一些历史无用的命令删除掉  
历史无用的命令如
```c
set name aa
set name bb //此时name aa 是无用的历史命令的
```
重写机制是，将某一时刻的内存数据，转换为命令，然后重新写到aof文件里

### RDB
即把某一时刻的内存状态，记录下来，持久化到磁盘里。  
持久化会另启动一种进程，所以持久并不会阻塞主线程  
持久化样例
```c
save 900 1 //900 秒之内，对数据库进行了至少 1 次修改；
save 300 10 //300 秒之内，对数据库进行了至少 10 次修改；
save 60 10000 //60 秒之内，对数据库进行了至少 10000 次修改。
```
### 混合持久化
为什么会有混合持久这种策略？
```
(1)RDB 优点是数据恢复速度快，但是快照的频率不好把握。频率太低，丢失的数据就会比较多，频率太高，就会影响性能。
(2)AOF 优点是丢失数据少，但是数据恢复不快。
```
混合持久
```
结合rdb+aof的优点:在aof重写日志时，以rdb方式记录下来，这段时间内主线程处理的写命令，会记录到一个缓冲区里，最后将缓冲区的的命令，再追加到文件里，即文件是rdb+aof
```

## redis集群
### 主从模式
```
优点：有数据备份安全、提高了读能力
缺点：master发生故障需手动切换
//注意在从节点，不会判断key是否过期，从节点过期key处理，依赖于主节点，主节点如果删除了过期key，会发送一条del指令到从节点
```
### 哨兵模式
![Alt text](image-3.png)
```
优点：在主从模式的基本上，实现了master故障自动切换，解决了高可用问题
缺点：性能及容量局现于单机
//注意：sentinel 是一个单独的进程，一般和redis同一台机器
```
#### 集群脑裂
```
当主机点与从节点网络有问题时，就会出现两个主节点的情况，即集群脑裂。
```
集群脑裂会导致什么问题？
```c
//min-slaves-to-write、min-slaves-max-lag
如果没有配置这两个参数或都这两个参数设置不合理，就会导致客户端向旧主节点写入数据丢失。
这两个参数的意义是控制主从同步数据的状态，如果有问题，就拒绝客户端写的请求。
min-slaves-to-write:主必须要有x个从连接点，如果少于拒绝写请求
min-slaves-max-lag：当主从复制数据超过x时间，就拒绝写请求
```
### cluster模式
![Alt text](image-2.png)
原理
```
redis集群分16384个槽点，每个集群节点负责一部分槽点；
槽点根据crc16(key)%16384得出
```
客户端
```c
客户端与redis直连，哪访问某个key不存在这个节点时，会返回 moved命令，告诉客户端正确的路由信息
//注意客户端会缓存集群槽的情况，当请求时也会算crc16(key)/16384。如果集群有变更，客户端可以发送 CLUSTER SLOTS 指令更新槽的情况
```
## redis过期删除策略
惰性删除+定期删除  
惰性删除
```
redis使用一个dict存储设置了过期时间的key，当客户端请求key的时候，判断key是否过期，如果过期了则删除key+返回null值
//优点：对cpu友好 
//缺点：对内存不友好
```
定期删除
```
(1)redis定期(默认10秒)随机从过期字典中随机抽取20个key，判断是否过期
(2)如果这次的抽取过期率超过25%那么继续下一次抽取
(3)如果随机抽取删除的操作耗时超过25ms，那会跳出本次检查
```
## 内存淘汰策略
是指如果redis使用的内存超过了上限，采取的措施。（注意redis 在64位默认不设置上限的）
```
可以分为三类：不处理，对设置过期key进行淘汰，对所有key进行淘汰
```
### 不处理
```c
redis3.0之后默认的策略 //noeviction
在此策略下，如果内存超过了上限，则所有写操作都会失败，读操作正常
```
### 设置过期key+所有key淘汰
这两种淘汰都差不多，[细节看](https://xiaolincoding.com/redis/module/strategy.html#%E5%86%85%E5%AD%98%E6%B7%98%E6%B1%B0%E7%AD%96%E7%95%A5)

### lru与lfu
lru
```
最近最少使用：按访问时间排序，时间越久的先被淘汰。
```
但lru算法有个bug
```
例如缓存大小为3，现有A,B,C三个元素，A和B被访问的次数都很多，突然C被访问了下，碰巧这时又要插入一个D，这时A或B会被淘汰掉。
```
#### redis 的lru算法
传统实现lru算法
```
需要维护一个有序链表，当元素被访问的时候，把元素移动表头，当要淘汰的时候就从表尾直接删除元素
```
redis的lru算法
```c
//传统的实现，有两个问题
(1)需要额外维护一个大链表
(2)每次访问需要移动元素
//redis实现类似lru算法
(1)元素增加一个字段保留最近访问时间
(2)当发现淘汰时，随机抽取N个元素，然后访问时间越久的，就淘汰掉
```
#### lfu
最近最不常用：按访问频率排序，频率越低的，先被淘汰掉

## redis设置key ttl总结
```
(1)内存淘汰策略，可以指定淘汰的范围只是设置了ttl的
(2)淘汰策略，支持lru和lfu等策略
(3)从节点不会主动删除或判断过期key，如果某个key过期了且主节点，没有发del命令过来，那么客户端读取从节点，是能获取到的
```
## redis缓存设计
### 应该避免 缓存雪崩、缓存击穿、缓存穿透这三个问题
#### 缓存雪崩
```
同一时间大量key过期，导致大量请求打到DB，把数据库压跨，进而系统连锁反应，最终导致系统不可用。
```
对应策略
```
过期时间设置为随机
不设置过期时间，后台线程定期更新
```
#### 缓存击穿
```
热点数据突然过期导致,大量请求打到DB
```
对应策略
```
(1)不设置过期时间，后台线程定期更新
(2)加个锁，只允许一个请求访问DB
```
#### 缓存穿透
```
即访问的数据即不在缓存，也不在数据库里
一般是非法用户或业务请求不合理 所造成
```
对应策略
```
(1)在api层限制非法ip
(2)当发现有缓存穿透时，设置一个返回默认值，不让他访问缓存和DB
(3)在写数据的时候，把key写到布隆过滤器里，查询数据库前，先判断数据在不在
```
### 如何保证缓存数据一致性
缓存不一致是指：DB的值和缓存致不一致(并且一般指不重启的情况，永久不一致)  
问题列表
```
(1)处理写请求时，数据写到数据库后，是更新缓存，还是删除缓存？
(2)处理写请求时，是先删除缓存，再写数据库，还是先写数据库，再删除缓存
(3)先更新数据库，再删除缓存，就能保证一致性吗
```
(1)问题分析
```c
//在两个并发写的情况下，更新缓存有脏数据问题，应该采取删除缓存
如果是采取更新缓存，那么在两个并发写的请求，那么可能缓存到旧值，所以应该采取删除缓存比较合适(get请求的时候，再load到缓存)
```
(2)问题分析
```c
//在并发读写的情况，先删除缓存，会有脏数据，应该采取先更新DB，再删除缓存
如果先删除缓存，再更新数据库，那么在：写请求将缓存删除了(未完成写到DB情况)，此时读请求发现缓存没中，从DB读(旧值)到缓存，此时写请求到DB才完成。
```
(3)问题分析
```c
//即使先更新数据库，再删除缓存，也还是有可能缓存不一致的。
例如:在删除缓存的时候失败(但实际这种情况概率比较小)

//妥协方案
缓存加个过期时间，这样即使不一致，也是某一段小时间

//稳妥方案
(1)加个删除缓存重试策略
(2)订阅mysql binlog日志，有消息时，更新对应的缓存
```

## redis 实战
### 如何实现延迟队列
延迟队列是指：如在淘宝、京东购物时，超过一段时间未付钱，订单自动取消
```
(1)用redis的有序集合，时间做为排序字段
(2)应用程序，根据当前时间来筛选，如果超时了，就进行删除业务
```
### bigkey  
bigkey的危害
```
(1)阻塞工作线程
(2)引发网络阻塞
(3)内存分配不均(cluster集群情况)
```
解决方案
```
//首先要制定使用规范
//做好监制
//删除bigkey时用unlink代替del命令(4.0以上)
//拆分思路
尽量拆分成多个key
```
### redis分布式锁
```
使用setnx实现，注意加个过期时间
value要用clientid，防止误删除
删除锁的时候，要用lua脚本，查询下对应的value，要和clientid相等时，才能删除
```



