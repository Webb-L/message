# 第一阶段：构建应用程序
FROM golang:1.21.7-alpine AS builder

# 设置作者信息
LABEL authors="webb"

# 设置工作目录
WORKDIR /app

# 复制所有文件到工作目录
COPY . .

# 设置环境变量
ENV GOPROXY=https://mirrors.cloud.tencent.com/go/
ENV GIN_MODE=release

# 构建应用程序
RUN go build -o message .

# 第二阶段：构建最终镜像
FROM alpine:latest

# 设置工作目录
WORKDIR /app

# 创建日志目录
RUN mkdir logs

# 从第一阶段复制配置文件、资源文件和应用程序到最终镜像
COPY --from=builder /app/config/config.yaml config/config.yaml
COPY --from=builder /app/resources resources
COPY --from=builder /app/message .

# 定义容器启动命令
CMD ["./message"]

# 暴露端口
EXPOSE 1204
