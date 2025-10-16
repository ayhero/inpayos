# 使用预构建的二进制文件
FROM alpine:3.17

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

# 复制预构建的二进制文件
COPY build/inpayos-linux ./inpayos

# 复制locales目录到镜像
COPY internal/locales ./internal/locales

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
CMD ["./inpayos"]