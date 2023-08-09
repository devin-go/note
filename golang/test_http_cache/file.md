# 分布式服务保障的三大手段
限流、熔断、服务降级
## 限流
限流的目的？
```go
//生活中例子：景区门票、医院挂号
保证有限的资源，能提供最大服务化的能力。在正常流量内的，提供服务，超过的则丢弃  
```
### 限流的算法
固定窗口  
```
限制单位时间内允许访问N次
如一秒内只允许访问10次
但这有临界值的问题，即1.9秒---2.9秒这一秒放行了20次。
```
滑动窗口
```
将单位时间划分为多个小周期，如将一秒划分为5个小周期，即0.2s一个周期，每过一个小周期就向前移，如1.9秒--2.9秒为一秒
解决了固定窗口的临界值问题
```
漏桶算法
```go
//所有请求都先到达一个队列里，然后服务再从队列里消费
ch:=make(interface{},10)//桶
func accept(req interface{}) {
    select{
        case ch<-req
        default
        //桶满了
    }
}
func handle() {
    tk:=time.NewTick(1s)
    for range tk {
        select{
            case req:=<-ch
            //handle req...
        }
    }
}
```
令牌算法
```
有一个池子，按照固定的速率往池子放令牌，然后理请求前，都来池子取令牌，取到令牌才处理，否则放弃这个请求。
```
![Alt text](image-4.png)

bbr 自适应算法
```
传统的算法都需要设置固定的限制值，但这个值通常双比较难取，bbr 算法通过系统的综合分析状态来判断是否限制流
系统状态一般包括：cpu、load、请求成功数，rt等(注意没有内存指标)
```

## 熔断
什么是熔断？  
![Alt text](image-3.png)
如图所示，当用户一个请求内部还要调用两个服务时才完成，但当调用服务C出现问题的次数到达一定量时，则服务A对它触发熔断，即不在一定的时间内不再调用服务C，从而保障整个系统(起码服务A和服务B是可用的)  

熔断的目的  
防止系统雪崩，如当服务C不可用是服务A还继续不停的访问，这样会导致服务A堆积了大量请求，进而把服务A的资源耗完。

### 熔断有三种状态  
```
关闭、打开、半打开
他们之间是循环关系
```
### 例子
```go
hystrix.ConfigureCommand("wuqq", hystrix.CommandConfig{
		Timeout:                int(3 * time.Second),//执行 command 的超时时间。
		MaxConcurrentRequests:  10,//command 的最大并发量 。
		SleepWindow:            5000,//当熔断器被打开后，SleepWindow 的时间就是控制过多久后去尝试服务是否可用了。
		RequestVolumeThreshold: 10,//一个统计窗口 10 秒内请求数量。达到这个请求数量后才去判断是否要开启熔断
		ErrorPercentThreshold:  30,//错误百分比，请求数量大于等于 RequestVolumeThreshold 并且错误率到达这个百分比后就会启动熔断
	})

	_ = hystrix.Do("wuqq", func() error {
		// talk to other services
		_, err := http.Get("https://www.baidu.com/")
		if err != nil {
			t.Logf("get error:%v",err)
			return err
		}
		return nil
	}, func(err error) error {
		t.Logf("handle  error:%v\n", err)
		return nil
	})
```

## 降级
```
从系统业务维度考虑，当流量增大时，可以关闭一些非核心服务非核心功能，从而保障核心服务或功能能正常服务。
也可以体现在，发生限流或熔断上。
```
