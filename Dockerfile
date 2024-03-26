# 这行设置基础镜像为golang:alpine，并将构建阶段命名为"builder"。golang:alpine镜像是一个轻量级镜像，包含了Go编程语言。
FROM golang:1.20.14-alpine AS builder

# 这行将默认的Alpine Linux仓库URL替换为Aliyun的镜像地址。这样做是为了提高中国用户的软件包下载速度。
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 这行更新Alpine的软件包索引并安装tzdata和bash软件包。tzdata用于设置时区，bash是一个常用的Shell解释器。ffmpeg是音视频处理要的
RUN apk update --no-cache && apk add --no-cache tzdata bash ffmpeg wget

# 这行设置工作目录为/workspace，并将当前目录中的文件复制到容器的/workspace目录中。
WORKDIR /workspace
COPY . .

# 这些行设置环境变量。GO111MODULE设置为 on，GOPROXY设置为https://goproxy.cn,direct以使用Go模块代理。
ENV GO111MODULE on
ENV GOPROXY https://goproxy.cn,direct

# 这行使用go mod命令下载依赖项模块。
RUN go mod download

RUN chmod +x build.sh && ./build.sh

#第二阶段
FROM ubuntu:20.04 AS production
WORKDIR /app

## 设置时区
RUN apt-get -y update && DEBIAN_FRONTEND="noninteractive" apt -y install tzdata
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

COPY ./*.sh  /app/

#本机目录不能使用绝对路径，因为它本身就是一个相对路径
#只会复制本机的config目录下所有文件，而不会创建config目录，所以后面需要指定
COPY ./config ./config
# COPY . .
RUN mkdir -p ./app/video/tmp/
RUN apt install -y wget ffmpeg gnupg systemctl

## 安装 etcd v3.5.9
RUN wget https://github.com/etcd-io/etcd/releases/download/v3.5.9/etcd-v3.5.9-linux-amd64.tar.gz &&\
    tar -zxvf etcd-v3.5.9-linux-amd64.tar.gz &&\
    cd etcd-v3.5.9-linux-amd64 &&\
    chmod +x etcd  &&\
    mv ./etcd* /usr/local/bin/

## 安装 Jaeger v3.5.9
RUN wget -c https://github.com/jaegertracing/jaeger/releases/download/v1.48.0/jaeger-1.48.0-linux-amd64.tar.gz &&\
    tar -zxvf jaeger-1.48.0-linux-amd64.tar.gz &&\
	cd jaeger-1.48.0-linux-amd64 &&\
    chmod a+x jaeger-* &&\
    mv ./jaeger-* /usr/local/bin/
    # nohup ./jaeger-all-in-one --collector.zipkin.host-port=:9411 &

## 安装 RabbitMQ
### 导入 RabbitMQ 的存储库密钥
RUN wget -O- https://github.com/rabbitmq/signing-keys/releases/download/2.0/rabbitmq-release-signing-key.asc | gpg --import -
### 将存储库添加到系统
RUN apt-get install -y apt-transport-https &&\
    cho "deb https://dl.bintray.com/rabbitmq-erlang/debian focal erlang" | tee /etc/apt/sources.list.d/bintray.erlang.list &&\
    echo "deb https://dl.bintray.com/rabbitmq/debian focal main" | tee /etc/apt/sources.list.d/bintray.rabbitmq.list
### 安装 RabbitMQ 和 Erlang
RUN apt-get install -y rabbitmq-server

## 安装 注意：Redis安装会自动启动
RUN apt install -y redis-server

COPY --from=builder /workspace/gateway .
COPY --from=builder /workspace/user .
COPY --from=builder /workspace/video .
COPY --from=builder /workspace/relation .
COPY --from=builder /workspace/favorite .
COPY --from=builder /workspace/comment .
COPY --from=builder /workspace/message .
EXPOSE 8080 16686

# RUN chmod +x /app/run.sh 等效下面语句
RUN cd /app &&chmod +x start.sh
CMD ["/app/start.sh"]

### docker build -t zhouxing9454/tiktok .
### docker run -it -p 8080:8080/tcp -p 16686:16686/tcp --name bytedance zhouxing9454/tiktok
