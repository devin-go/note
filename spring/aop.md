# aop
本质就是代理模式，例如
拦截类
```java
@Aspect
@Component
public class LoggingAspect {
    // 在执行MailService的每个方法前后执行:
    @Around("execution(public * com.itranswarp.learnjava.service.MailService.*(..))")
    public Object doLogging(ProceedingJoinPoint pjp) throws Throwable {
        System.err.println("[Around] start " + pjp.getSignature());
        Object retVal = pjp.proceed();
        System.err.println("[Around] done " + pjp.getSignature());
        return retVal;
    }
}
```
被代理类
```java
@Component
@EnableAspectJAutoProxy
public class MailService {
    public void sendMsg(msg string) {

    }
}
```
模拟spring自动生的代理类
```java
public class AutoProxyMailServicexxx extends MailService{
    MailService *target;
    LoggingAspect *aspect;

    public void sendMsg(msg string){
        aspect.doLogging(target)//实际是wrapTarget对象
    }
}
```
## 自定义拦截注释标签
上面的例子，是根据表达式来拦截(拦截了MailService类下的所有方法)，这样有时候会不太清析，如后面要方法，但不需要aop拦截。
解决方法是，可以自己实现拦截注释标签，参考例子是spring 的事务标签
```java
@Component
public class UserService {
    // 有事务:
    @Transactional
    public User createUser(String name) {
        ...
    }
    // 无事务:
    public boolean isValidName(String name) {
        ...
    }
}
或者直接定义到类上(这个类的所有方法都有事务了)
@Component
@Transactional
public class UserService {
    public User createUser(String name) {
        ...
    }
    public boolean isValidName(String name) {
        ...
    }
}
```

## aop拦截器类型
```
before //执行前
after//执行后
afterReturn//执行正确后执行(即未抛异常)
afterThrow//执行有抛异常后执行
around //环绕，即包含上面所有的功能
```
## aop避坑点
```
1.final 类 不能被代理，所以要想使用aop的业务类，不能定义为final类
2.不要使用字段
```
### 不要使用字段，模拟代码
错误的
```java
public class A {
    Object o = NewObject();
}
public class B {
    @Autowired
    A a;
    public void xxx() {
        r:=a.o.xxx();//如果被aop代理了，那么o就为空对象
    }
}
```
正确的
```java
public class A {
    Object o = NewObject();
    public Object getO() {
        return o;
    }
}
public class B {
    @Autowired
    A a;
    public void xxx() {
        r:=a.getO().xxx();//如果被aop代理了，那么o就为空对象
    }
}
```
[详细原因](https://www.liaoxuefeng.com/wiki/1252599548343744/1339039378571298#0)