# sync.pool
```
(1)利用 GMP 的特性，为每个 P 创建了一个本地对象池 poolLocal，尽量减少并发冲突
(2)localPool有一个链表(poolChain)，链表元素是一个ringBuffer
(3)put对象时，如果ringBuffer满了，则会创建一个新的
(4)get流程，先获取本地的-->窃取其他P的-->上一轮的缓存池获取-->调用户New方法
```
## poolChain示意图
![Alt text](image-5.png)
## get流程
![Alt text](image-6.png)

# sync.Mutex
```
加锁的原理是cas设置一个变量值为1，设置成功则获取到锁
获取不到锁的协程，则进入等待队列
锁有两种模式：正常模式和饥饿模式
```
正常模式
```
当释放锁的时候，会从等待队列唤醒一个G，但这个G可能跟后面来的G抢锁
而后来的G，可能处于自旋状态，更有机率获取到锁
```
饥饿模式
```
当从等待队列唤醒G，排队时间超过1MS且没获得到锁，就会将锁标记为“饥饿模式”
那么后面来的G就不自旋和参与抢锁，而是直接插入到队列尾部
```
饥饿模式-->正常模式
```
当从等待队列唤醒的G等待时间没超过1MS时，就将锁切回正常模式
当等待队列为空后，也会将锁切回正常模式
```
# sync.WaitGroup
原理易于理解版
```go 
type WaitGroup struct {
  waiter int
  counter int
  sema int
}
调用wg.Add(n)时 相当于counter+=n
调用wg.Done()时 相当于counter--，如果counter==0，会释放sema个信号量唤醒调用wg.Wait()的Goroutine
调用wg.Wait()时 相当于waiter++，同时增加信号量，并挂起当前goroutine
```
无锁细节
```
试想如果想两个变量支持并发，直接用锁就行了，但用锁会增加性能开销，golang的做法是：
  counter和waiter柔和到一个int64变量里，高32代表counter，低32代表waiter
```
注意
```go
Add和Wait方法，不能同时调用，否则会panic//自已看Add源码看到
```
# slice
扩容机制
```
初始长度         增长比例
256             2.0
512             1.63
1024            1.44
2048            1.35
4096            1.30
```
解释：小于256则翻倍，其他长度扩容类似
# map
特点（[原理](https://qcrao91.gitbook.io/go/map/map-de-di-ceng-shi-xian-yuan-li-shi-shi-mo)）
```c
(1)hash(key)后，用低N位bit确定桶的位置
(2)桶有个tophash数组([8]byte)用来存储高8位的hash值(这样的好处是提高检索key的效率)
(3)每个桶固定能存储8个kv
(4)桶内是按照/k/k/.../v/v/..来存储的，并不是按k/v这样存储//golang解析说，这样能有效节省空间(内存对齐)
(5)扩容方式是渐近式的(类似redis的dict)
//如何计算map有多少个桶？
如果N是5，那么有2^5=32个桶
//tophash 为什么能提高检索key的效率
(1)经过hash(key)定位到具体桶后，key的前辍大概率是不一样的
(2)整形比较给其他类型快
```
key定位过程
```
(1)根据hash(key)得到64位的二进制
(2)低N位用来定位到桶
(3)桶固定有8个位置，还有一个固定长度为topbits的整形数组，
(3.1)topbits存放高8位的值，查找时先比较topbit，如果相等，则再进一步比较key(比较前辍key思想)
(3.2)如果不等，则继续找下一位
(3.3)如果topbits里找不到且溢桶不为空，也要在溢出桶里找
```
key冲突处理
```
(1)如果桶还有空的位置，则存放到桶里
(2)如果桶不为空，则存放到溢出桶
```
delete过程
```
根据key定位到位置，然后将key和value“清零”
```
赋值过程
```

```
扩容
```c
//扩容条件
(1)负载因子超过6.5，count/2^B>=6.5，翻倍扩容
(2)溢出桶过多(新增元素后，又大量删除)，等量扩容。
(2.1)如果B小于15，溢出桶大于2^B,则也会触发扩容
(2.2)如果B大于等于15时，溢出桶超过2^15也会触发扩容
//扩容过程
(1)标记map正在扩容
(2)采用渐进式迁移,一次最多迁移两个桶(注意迁移桶时，也会迁移溢出桶)
(3)如果是翻倍扩容(即相对于原来多了一位bit)，有个小技巧，多一位bit如果是0，则新的桶编号和原来的编号一样，如果是1，是翻一位后对应
```
## map 一些注意事项
同一个协程，可以边遍历边删除吗？
```go
for k,v:=range m{
  delete(m,k)
}
答：是可以的
原理：删除只是将key和value“清零”，底层数组未变化
```
为什么每次遍历map都是无序的？(即使这个map只是只读)
```c
因为每次遍历都会从一个随机桶+随机cel开始。
//为什么要随机？
因为 map存在扩容的情况，且会有中间状态
在中间状态下，每次遍历都会不一样
所以为了兼容这种表现，go每次遍历都会随机开始
```
float64可以做为map的key吗
```c
答：可以的
原因：如果key是float64，map会将float64转为unint64（math.Float64bits()）
      两个float64，如math.Float64bits(2.4)和math.Float64bits(2.4000000000000000000000001)会相等
      所以慎用float64作为key
//math.NaN情况
m[math.NaN]=1,m[math.NaN]=2，map会存两个
原因是map判断是NaN时，会加一个随机值，hash(randNum)
//key是必须是可以比较的
如slice,map 是不可以做为key的
```
# sync.Map
```c
一般情况实现并发安全的map+mutex，但这样读和写都要加锁，效率有点低。
sync.Map 内部维护read和dirty两个map，实现了读写分离
//读操作,对应Load方法
(1)优先从read取 no lock
(2)如果从read没取到，再从dirty里取 have lock
//写操作，对应strore方法，有三种情况
(1)如果read里有，则CAS更新，no lock
(2)如果dirty里有，则更新，have lock
(3)如果read和dirty都没有：have lock
    将read的值拷贝到dirty
    将新值存入到dirty，并标记read和dirty不相等
//删除操作
(1)如果read里有，标记为删除
(2)如果dirty里有，直接删除

//注意
Load方法会统计从dirty取的次数，当超过一定次数时，会将dirty赋值给read
```
优缺点
```
适合读多写少场景
```
# context
```
对外有六种context：
  cancel,deadline,timeout,value,todo,backgroud
内部有四种
  emptyCtx,valueCtx,cancelCtx,timerCtx
其中 timeCtx继承了cancelCtx
timout,deadline 对应了 timerCtx
todo,backgroud对应了emptyCtx
```
timeout和deadline context原理
```
使用了 time.After()关闭done chanl
```
父context cancel，子context也跟着取消
```
首先new子context时传入了父的ctx,如context.Withxx(parentCtx)，然后会启动一个协和监听父是否cancel了，如果父cancel了，自己也会cancel掉
```
多个context间同样的key，value不会覆盖，但获取时优先取到最近的
```
根据key获取value时，类似链表从送开始往后查找，当查到了就返回了
```
# channel
## 原理
结构
```go
type hchan struct {
  lock mutex //锁
  buf      unsafe.Pointer//循环链表
  sendx    uint   // 发送 buf的位置
  recvx    uint   // 接收 buf的位置
  recvq    waitq  // 接收gorutine队列
  sendq    waitq  // 发送gorutine队列
}
```
发送流程
```
(1)加锁
(2)将值拷贝到buf
(3)如果buf满了，会阻塞，并加入sendq队列
(4)从recvq唤醒一个G
```
接收流程
```
(1)加锁
(2)从buf里拷贝值
(3)如果buf为空，会阻塞，并加入recvq队列
(4)从sendq唤醒一个G
```
close
```
(1)唤醒所有读队列recvq的G
(2)唤醒所有写队列sendq的G，并且写的G会panic
```
## channel三种状态和三种操作结果
| 操作     | 空值(nil) | 已关闭   | 未关闭 |
| -------- | --------- | -------- | ------ |
| 关闭     | panic     | panic    | 正常   |
| 发送数据 | 永久阻塞  | panic    | 正常   |
| 读取数据 | 永久阻塞  | 永不阻塞 | 正常   |

# interface
```
interface变量内部包含了两个字段：类型T与值V
当两个interface比较时，先比较T再比较V
当“其他变量”与interface比较时，先将“其他变量”转换为interface再比较
```
当一个类型指针赋给interace后，判断是否为nil
```go
  var p *int
	var i interface{}
	i = p
	t.Log(i==nil)//false
//为什么是false?
因为首先要将nil转换为interface即(T=nil,V=nil)，而i的结构为(T=*int,V=nil)，所以不相等
```
# struct{}
空struct,最主要只是节省内存而已
```c
//为什么struct{}占用空间为0
所有空struct 都指向：runtime.zerobase 
//struct{}的作用
 map做set
 chan 信号
```
# go语言初始化顺序
```
import-->const-->var-->init-->main

同包init顺序，不保证
```
# golang 有两种字符类型
```
英文字符类型用1个byte（即一个字节就可以表示 ASCII）
其他字符(1~4字节)，如中文(3字节)要用rune
```
# defer
```
在函数返回前运行，触发时机有三个：
（1）遇到return
（2）函数末尾
（3）遇到panic
defer对应有一个_defer结构体，实现defer有三个版本：
go13之前，_defer结构体会在堆上分配，执行时需要将堆上的变量拷贝到栈上，性能比较差
go13 _defer结构体，实现了在栈上分配，性能相比之前提升了30%
go14实现了开放式编码(直接将derfer 函数在当前函数展开)，性能损失几乎可以不计，但也有条件：
  (1)未关必内关联
  (2)没在for循环里调用defer
  (3)函数乘积没超过15
```
# panic
对应了runtime._panic结构体
```go
type _panic struct {
	argp      unsafe.Pointer //指向defer 指针
	arg       interface{}// 调用panic时的参数
	link      *_panic//上一个panic
	recovered bool
	aborted   bool
	pc        uintptr
	sp        unsafe.Pointer
	goexit    bool
}
```
嵌套panic
```go
func main() {
	defer fmt.Println("in main")
	defer func() {
		defer func() {
			panic("panic again and again")
		}()
		panic("panic again")
	}()
	panic("panic once")
}
输出结果：
panic once
panic again
panic again and again
in main
```

## defer,panic,recover
如果没有recover,panic 会导致主程序退出。并且recover只能捕获当前goroutine
```go 
func f() {
  defer print("f defer")
  panic("aa")
}

func main(){
  defer print("main defer")
  go f()
  time.sleep(1)
}
//结果不会输出 main defer
```
recover只在panic之后生效
```go
func main(){
  if err:=recover();err != nil {
    print(err)
  }
  panic("aa")
}
//结果recover捕获不到err
正确的例子：
func main() {
  defer func(){
    if err:=recover();err!= nil {
      print(err)
    }
  }
  panic("aa")
}
```
# ==,比较
一个比较相通的原则是，[比较详细的介绍](https://darjun.github.io/2019/08/20/golang-equal/)
```
只有类型相同且值相同时，==比较才是true
```
golang可以分为四大类型：
```
基本类型：int,float,string
复合类型：array和 struct
引用类型：pointer，slice,map,channel
接口类型：如error
```
注意
```
两个slice之间不能比较，只能和nil比较
两个map之间不能比较，只能和nil比较
```
两个interface比较
```
(1)类型要相等
(2)值要相等(用的是==比较法)
(3)如果interface指向的是slice或map会报panic，因为两个slice或两个map之间是不可以比较的
```
# 什么是协程
```
协程的执行是在线程上的，
一个线程里可以有多个协程，
线程内的协程是串行的执行的
简单可以理解为用户态的线程。
```
协程切换
```
协程切换的代价比较小，可以简单理解为只切换了寄存器和协程栈的内容
```

# select 原理
[原理](https://cloud.tencent.com/developer/article/2205909)
```go
type scase struct{
  c channel
  ...
}
select 语句用于channel上，用于监听多个channel是否可写或可读，思想和io多路复用一样
case对应一个结构体runtime.scase
如果有多个case，最终会运行到runtime.selectgo()函数上
  该函数会随机开始遍历case数组，保证channel不会有饿死的情况
  同时也会对根据channel的地址加锁，保证不会产生死锁
```
# 一些知识
```
变量是分配到堆上还是栈上，是由编译阶段决定的，golang编译会做逃逸分析
  如函数返回指针，一般会分配到堆上
栈
  一般指函数栈
  函数栈的作用：保存局部变量、向被调函数传递参数、保存函数返回地址等
  栈的大小随着函数调用层级增加而增大，随着函数的返回而减小。
  函数栈自动回收的关键是：编译器插入了代码实现的(c++函数栈也会自动回收)
go函数返回指针安全吗
  安全的，因为值分配在堆上
```
# cas 原理
```
cas = CompareAndSwap(v,e,n),依靠硬件实现：
    只有v=e才会将值改为n，否则什么都不做

但cas有ABA问题：
  即有三个线程同时修改一个变量，线程1将值变为A-->B，线程3将值变为B-->A，线程2将值变为A-->B
  这时线程2是不知道，变量中间有变化过
  解决方法是：给变量增加版本号，参考java的AtomicStampedReference实现
```

atomic.Addint(&v,delta)原理
```go
  利用cas实现
  func Addint(int *addr, delta) {
      for {
        expect:=*addr
        v = *addr + delta
        if CompareAndSwap(addr, expect, v) {
          break
        }
      }
  }
```
# sync.Value原理
前提
```go
interface{}变量，go内部会转换成类似这个结构体
  type efaceWords struct {
    typ  unsafe.Pointer
    data unsafe.Pointer
  }
typ 的作用是中间状态
```
Strore()原理
```
  第一次的情况：
    如果是没赋过值，即typ==nil，则会用CompareAndSwap将typ设置为“中间临时状态”
    然后再用StorePointer()赋值data和type
    (其他线程判断typ==中间临时状态，会重新来，即等待的意思)
  后面的情况：都是调用StorePointer()赋值data
```
Load()原理
```
LoadPointer()
```
Swap()原理
```
SwapPointer()
```