# Tiktok
自己独立修改的抖声项目，基于青训营团队抖声项目。原项目是我们青训营团队的大项目作业，这个是原项目地址：https://github.com/Lionel24-xxy/douyin-project





### 修改1

将sha-1改为sha-256算法，加盐值，盐值和password一起存入数据库。

```Go
func HashPassword(password string, salt string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password + salt))
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)
}

func GenerateSalt(length int) string {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(salt)
}
```





### 修改2

jwt的库修改（原来的库已经不再维护了）

```go
import "github.com/golang-jwt/jwt"
```





### 修改3

删除normal.go，将对应的中间件放入AuthMiddlerware里面





### 修改4

优雅地重启或停止web服务器，参考gin的官方文档写的：

- 启动了一个 goroutine，用于监听和处理客户端请求。如果在监听过程中出现错误，且该错误不是 `http.ErrServerClosed`，则会使用 `log.Fatalf` 记录错误信息并退出程序。接着，代码使用 `signal` 包等待中断信号。一旦接收到中断信号，代码会优雅地关闭服务器，首先使用 `context` 包创建一个超时为 5 秒的上下文，然后调用 `srv.Shutdown` 方法，等待服务器处理完尚未处理完的请求并关闭服务器。最后，程序记录一条日志消息表明服务器已关闭，并退出程序。





### 修改5

将使用计数器限流的方式改为使用令牌桶限流。

我们有一个固定的桶，桶里存放着令牌（token）。一开始桶是空的，系统按固定的时间（rate）往桶里添加令牌，直到桶里的令牌数满，多余的请求会被丢弃。当请求来的时候，从桶里移除一个令牌，如果桶是空的则拒绝请求或者阻塞。

令牌桶有以下特点：

- 令牌按固定的速率被放入令牌桶中
- 桶中最多存放 B 个令牌，当桶满时，新添加的令牌被丢弃或拒绝
- 如果桶中的令牌不足 N 个，则不会删除令牌，且请求将被限流（丢弃或阻塞等待）

**令牌桶限制的是平均流入速率**（允许突发请求，只要有令牌就可以处理，支持一次拿3个令牌，4个令牌...），**并允许一定程度突发流量，所以也是非常常用的限流算法。**

效果：

![image-20230423212527214](https://blog-1314857283.cos.ap-shanghai.myqcloud.com/images/202304232125307.png)
