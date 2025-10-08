# 多阶段构建
FROM golang:1.23-alpine AS builder

# 安装git和ca-certificates
RUN apk add --no-cache git ca-certificates

# 设置工作目录
WORKDIR /app

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN mkdir -p build && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/inpayos ./main

# 最终阶段
FROM alpine:latest

# 安装ca-certificates和时区数据
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 创建非root用户
RUN addgroup -g 1001 appgroup && adduser -u 1001 -G appgroup -s /bin/sh -D appuser

# 创建应用目录和日志目录
RUN mkdir -p /app/logs /logs

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/build/inpayos ./inpayos

# 复制locales目录到镜像
COPY --from=builder /app/internal/locales ./internal/locales

# 创建日志目录符号链接，支持两种挂载方式
RUN ln -sf /logs /app/logs-external || true

# 设置权限
RUN chown -R appuser:appgroup /app && chown -R appuser:appgroup /logs

# 切换到非root用户
USER appuser

# 暴露端口 (OpenAPI: 8080, CashierAPI: 8081, MerchantAPI: 8082, CashierAdminAPI: 8083, AdminAPI: 8084)
EXPOSE 8080 8081 8082 8083 8084

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 运行应用
CMD ["./inpayos", "serve"]