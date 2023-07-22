# 什么叫ioc
## 背景
```
如项目有 n api接口，每个api接口 都有可能依赖相同的m service 共同对象，而servcie 对象同样也可能依赖z 个 共同dao对象
```
传统的解决方法
```
需要手动new 每一个对象，如在 a api 里new a service 对象。
手动new 还可能会产生不必要的重复对象，如单例场景(虽然可能通过写代码解决，只实例单例，但麻烦)
```
## spring ioc
解决方法
```
通过配置文件(xml) 或 注解 注入要依赖的对象
实际项目中，通过注解比较多
```
支持场景
```
单例:对每次注入的都是同一个对象
多例:即每次注入的都是不同的对象
按条件注入
```
### 按条件注入
#### profile
[用来解决区分不同环境，使初化加载不同的配置](https://juejin.cn/post/6844904128238501896)

实际中建议不要用这种方式，侵入性太大了
好的解决方法是：将配置写到etcd里或文件里，不同的环境配置值不一样

#### conditional
根据条件进行初始化实例
```
如有一个上传文件的功能，本地环境时，就存储到本地目录下，生产环境就存储到s3等
```
```java
//上传文件接口
public interface Uploader{
    public void Upload(xxx) throws Exception;
}
//本地实现
@Component
@ConditionalOnProperty(name = "app.storage", havingValue = "file", matchIfMissing = true)
public class LocalStroe implements Uploader
//s3实现
@Component
@ConditionalOnProperty(name = "app.storage", havingValue = "s3")
public class S3Stroe implements Uploader
//使用uploader
@Component
public class UserImageService {
    @Autowired
    Uploader uploader;
}
```
ConditionalOnProperty其实是实现了spring的一个Condition 接口来实现的，我们也可以自定义化
## 总结
```
ioc 其实就是解决了各种bean 在各层之间的实例化的问题，支持的场景有：单例，多例，按条件初始化
golang 这边暂时需要手动new
```
