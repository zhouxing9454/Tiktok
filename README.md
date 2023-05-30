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





### 修改6

使用Docker容器化技术，将项目部署到Docker里面

**docker-compose.yml**

```yml
version: "3.9" # 使用3.9版本的docker-compose文件格式
services: # 定义服务
  web: # web服务
    build: .
    ports: # 映射端口
      - "8000:8000"
    volumes:
      - ./static/:/app/server/static/
    depends_on: # 定义web服务依赖的其他服务
      mysql_a: # web服务依赖mysql服务
        condition: service_healthy # web服务只有在mysql服务的healthcheck状态为healthy时才启动
      redis_a: # web服务依赖redis服务
        condition: service_started # web服务只有在redis服务启动后才启动
  redis_a: # redis服务
    image: "redis:alpine" # 使用redis:alpine镜像
    ports: # 映射端口
      - "6379:6379"
    restart: always # 总是重启
    environment: # 定义环境变量
      - REDIS_PASSWORD=zx045498  # redis密码
    # 这行指定了要在容器启动时执行的命令。具体地，它运行redis-server命令，并通过--requirepass选项设置Redis服务器的密码为zx045498。
    #这样，在启动该容器时，Redis服务器将以指定的密码进行身份验证。这有助于保护Redis实例免受未经授权的访问
    command: redis-server --requirepass zx045498
  mysql_a: # mysql服务
    image: "mysql:latest" # 使用mysql:latest镜像
    ports: # 映射端口
      - "3307:3306"
    restart: always # 总是重启
    environment: # 定义环境变量
      - MYSQL_ROOT_PASSWORD=123456 # 设置mysql的root密码为123456
      - MYSQL_DATABASE=byte_dance # 设置mysql启动后会默认创建一个byte_dance的database
      - MYSQL_ROOT_HOST=% # 设置任何机器都可以连接当前数据库
    healthcheck: # 定义mysql服务的健康检查
      test: [ "CMD", "mysql", "--user=root", "--password=123456", "--execute", "SHOW DATABASES;" ] # 使用mysql命令来检查数据库是否可用
      interval: 10s # 每10秒执行一次检查
      timeout: 5s # 检查超时时间为5秒
      retries: 3 # 检查失败后重试3次
```



**Dockerfile**

```shell
# 这行设置基础镜像为golang:alpine，并将构建阶段命名为"builder"。golang:alpine镜像是一个轻量级镜像，包含了Go编程语言。
FROM golang:alpine AS builder

# 这行为构建阶段添加一个标签，指示它是"gobuilder"阶段。
LABEL stage=gobuilder

# 这些行设置环境变量。CGO_ENABLED设置为0以禁用CGO（C Go）支持，GOPROXY设置为https://goproxy.cn,direct以使用Go模块代理。
ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn,direct

# 这行将默认的Alpine Linux仓库URL替换为Aliyun的镜像地址。这样做是为了提高中国用户的软件包下载速度。
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 这行更新Alpine的软件包索引并安装tzdata和bash软件包。tzdata用于设置时区，bash是一个常用的Shell解释器。
RUN apk update --no-cache && apk add --no-cache tzdata bash

# 这行设置工作目录为/build，并将当前目录中的文件复制到容器的/build目录中。
WORKDIR /build
COPY ../../ .

# 这行使用go mod命令下载依赖项模块。
RUN go mod download

# 这行使用go build命令编译Go应用程序，并使用-ldflags="-s -w"参数设置链接标志，以进行静态编译。编译后的二进制文件将被复制到/app/server路径下。
RUN go build -ldflags="-s -w" -o /app/server


# 这行设置基础镜像为alpine。
FROM alpine

# 这些行从之前的构建阶段中复制证书和时区信息到当前镜像中，并通过设置TZ环境变量将时区设置为"Asia/Shanghai"。
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai

# 这行指定容器将监听的端口号为8000。（暴露端口）
EXPOSE 8000

#  这些行设置工作目录为/app，并从之前的构建阶段中复制敏感词文件、应用程序二进制文件以及静态文件到容器的相应路径中。
WORKDIR /app
COPY --from=builder /build/utils/sensitiveWords.txt /app/server/utils/sensitiveWords.txt
COPY --from=builder /app/server /app/webserver
COPY --from=builder /build/static/ /app/server/static/

# 这两行安装了curl和ffmpeg软件包，使用--no-cache选项表示不缓存安装包。
RUN apk add --no-cache curl
RUN apk add --no-cache ffmpeg

# 这行设置容器启动时执行的命令，即运行webserver可执行文件。
CMD ["./webserver"]
```

