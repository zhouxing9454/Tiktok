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