FROM golang:latest AS builder
LABEL org.opencontainers.image.author="turbo514"
LABEL org.opencontainers.image.created="2025.11.13.v1"

WORKDIR /build/tenant

# 复制go模块文件
COPY internal/tenant/go.mod internal/tenant/go.sum ./
COPY internal/shared/go.mod internal/shared/go.sum ../shared/

# 下载依赖
RUN GOPROXY=https://goproxy.cn,direct go mod download

# 复制所有源代码
COPY internal/tenant/ ./
COPY internal/shared/ ../shared/

# 构建应用程序
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/app ./main.go

FROM alpine:latest
WORKDIR /app

# 复制二进制文件
COPY --from=builder /app/app .

# 复制配置文件
#COPY internal/tenant/config/config_prod.yaml ./config/config.yaml
#COPY internal/shared/commonconfig/global_prod.yaml ./config/global.yaml

EXPOSE 30002 30003

CMD ["./app"]