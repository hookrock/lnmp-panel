# 多阶段构建Dockerfile - ARMv7l优化

# 构建阶段
FROM golang:1.19-alpine AS builder

# 设置ARMv7l构建参数
ARG TARGETARCH=arm
ARG TARGETVARIANT=v7
ARG GOARM=7

WORKDIR /app

# 复制go模块文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建ARMv7l优化版本
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 \
    go build -a -installsuffix cgo -o lnmp-panel main.go

# 运行阶段
FROM alpine:latest

# 安装必要的系统工具
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 创建运行用户
RUN addgroup -S lnmp && adduser -S lnmp -G lnmp

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/lnmp-panel .

# 创建必要的目录
RUN mkdir -p /var/log/lnmp-panel

# 切换用户
USER lnmp

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 启动命令
CMD ["./lnmp-panel"]