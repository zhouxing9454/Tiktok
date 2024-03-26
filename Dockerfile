FROM golang:1.20.14-alpine AS builder

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk update --no-cache && apk add --no-cache tzdata bash ffmpeg wget

WORKDIR /workspace
COPY . .

ENV GO111MODULE on
ENV GOPROXY https://goproxy.cn,direct

RUN go mod download

RUN chmod +x build.sh && ./build.sh

FROM ubuntu:20.04 AS production

WORKDIR /app

RUN apt-get -y update && DEBIAN_FRONTEND="noninteractive" apt -y install tzdata
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

COPY ./*.sh  /app/

COPY ./config ./config

RUN mkdir -p ./app/video/tmp/

RUN apt install -y wget ffmpeg gnupg systemctl

RUN wget https://github.com/etcd-io/etcd/releases/download/v3.5.9/etcd-v3.5.9-linux-amd64.tar.gz &&\
    tar -zxvf etcd-v3.5.9-linux-amd64.tar.gz &&\
    cd etcd-v3.5.9-linux-amd64 &&\
    chmod +x etcd  &&\
    mv ./etcd* /usr/local/bin/

RUN wget -c https://github.com/jaegertracing/jaeger/releases/download/v1.48.0/jaeger-1.48.0-linux-amd64.tar.gz &&\
    tar -zxvf jaeger-1.48.0-linux-amd64.tar.gz &&\
	cd jaeger-1.48.0-linux-amd64 &&\
    chmod a+x jaeger-* &&\
    mv ./jaeger-* /usr/local/bin/

RUN wget -O- https://github.com/rabbitmq/signing-keys/releases/download/2.0/rabbitmq-release-signing-key.asc | gpg --import -

RUN apt-get install -y apt-transport-https &&\
    cho "deb https://dl.bintray.com/rabbitmq-erlang/debian focal erlang" | tee /etc/apt/sources.list.d/bintray.erlang.list &&\
    echo "deb https://dl.bintray.com/rabbitmq/debian focal main" | tee /etc/apt/sources.list.d/bintray.rabbitmq.list

RUN apt-get install -y rabbitmq-server

RUN apt install -y redis-server

COPY --from=builder /workspace/gateway .
COPY --from=builder /workspace/user .
COPY --from=builder /workspace/video .
COPY --from=builder /workspace/relation .
COPY --from=builder /workspace/favorite .
COPY --from=builder /workspace/comment .
COPY --from=builder /workspace/message .

EXPOSE 8080 16686

RUN cd /app &&chmod +x start.sh
CMD ["/app/start.sh"]

### docker build -t zhouxing9454/tiktok .
### docker run -it -p 8080:8080/tcp -p 16686:16686/tcp --name bytedance zhouxing9454/tiktok

#这段Dockerfile是用于构建一个基于Golang和Ubuntu的容器化应用。以下是每个命令的中文解释：
 #
 #1. `FROM golang:1.20.14-alpine AS builder`：
 #   - 使用基础镜像 `golang:1.20.14-alpine`，并将此构建阶段命名为 "builder"。
 #   - Alpine镜像是轻量级的，并包含了Go编程语言。
 #
 #2. `RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories`：
 #   - 替换默认的Alpine Linux仓库URL为阿里云的镜像地址，以提高中国用户的软件包下载速度。
 #
 #3. `RUN apk update --no-cache && apk add --no-cache tzdata bash ffmpeg wget`：
 #   - 更新Alpine的软件包索引并安装 `tzdata`、`bash`、`ffmpeg` 和 `wget` 软件包。
 #   - `tzdata` 用于设置时区，`bash` 是一个常用的Shell解释器，`ffmpeg` 是音视频处理要的。
 #
 #4. `WORKDIR /workspace` 和 `COPY . .`：
 #   - 设置工作目录为 `/workspace`，并将当前目录中的文件复制到容器的 `/workspace` 目录中。
 #
 #5. 设置环境变量：
 #   - `ENV GO111MODULE on`：启用Go模块支持。
 #   - `ENV GOPROXY https://goproxy.cn,direct`：设置Go模块代理。
 #
 #6. `RUN go mod download`：
 #   - 使用go mod命令下载依赖项模块。
 #
 #7. `RUN chmod +x build.sh && ./build.sh`：
 #   - 修改 `build.sh` 文件为可执行，并执行该脚本。这里假设 `build.sh` 负责构建你的Go应用程序。
 #
 #8. 第二阶段 `FROM ubuntu:20.04 AS production`：
 #   - 使用基础镜像 `ubuntu:20.04`，并将此阶段命名为 "production"。
 #
 #9. `WORKDIR /app`：
 #   - 设置工作目录为 `/app`。
 #
 #10. 设置时区：
 #    - 更新软件包索引并安装 `tzdata` 软件包，设置系统时区为亚洲/上海。
 #
 #11. 复制脚本和配置文件：
 #    - 将 `.sh` 文件复制到 `/app` 目录。
 #    - 复制 `config` 目录下的文件到容器的 `./config` 目录下。
 #
 #12. 安装依赖：
 #    - 安装 `wget`、`ffmpeg`、`gnupg` 和 `systemctl`。
 #    - 下载并安装 `etcd v3.5.9` 和 `Jaeger v3.5.9`。
 #    - 添加并安装 `RabbitMQ` 和 `Erlang`。
 #    - 安装 `Redis`。
 #
 #13. `COPY --from=builder ...`：
 #    - 从前一个构建阶段复制编译好的Go应用程序到 `/app` 目录。
 #
 #14. `EXPOSE 8080 16686`：
 #    - 暴露容器的端口 `8080` 和 `16686`。
 #
 #15. `RUN cd /app && chmod +x start.sh`：
 #    - 进入 `/app` 目录并修改 `start.sh` 文件为可执行。
 #
 #16. `CMD ["/app/start.sh"]`：
 #    - 容器启动时执行 `start.sh` 脚本。