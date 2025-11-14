# 项目简介
这是一个 SaaS 多租户短链接平台。每个租户可创建和管理短链接，终端用户点击短链接后系统会进行 跳转、限流、事件上报、监控日志记录 等操作
该项目用于学习微服务、Golang、高并发、缓存策略等相关技术

# 功能特性
- 多租户管理：独立空间、独立短链、独立统计数据
- 短链生成 & 跳转服务
- 访问日志 & 点击事件上报
- Redis 缓存加速 / 缓存穿透保护
- 本地缓存加速
- Prometheus + Grafana 指标监控
- OpenTelemetry 分布式追踪
- 数据库迁移工具
- Redis实现令牌桶算法分布式限流
- cqrs
- clean architecture

# 技术栈
## Backend
- Golang（Gin / gRPC）
- GORM
- Viper

## Infra
- Redis
- MySQL
- RabbitMQ
- Clickhouse

## DevOps
- Docker & Docker Compose
- Makefile
- golang-migrate

## Observability
- Prometheus
- OpenTelemetry + Jaeger

# 性能表现

可以做到600qps

由于是本机测试,应该有一定误差
一大部分性能损耗来自消息队列和序列化

几乎全部命中缓存

# 快速开始
```bash
git clone https://github.com/turbo514/shortenurl-v2.git

cd shortenurl-v2

make prod-infra-up

make prod-app-up
```

```bash
# 创建租户
curl -X POST "http://localhost:8080/tenants" \
    -H "Content-Type: application/json" \
    -d '{"name":"张三"}'

# 创建用户
curl -X POST "http://localhost:8080/users" \
  -H "Content-Type: application/json" \
  -d '{
        "name":"张三",
        "tenant_id":"019a829b-2231-7740-ba47-243aca9eadeb",
        "api_key":"74bf9cf5-c310-4a48-b904-06cc4357c933",
        "password":"123456"
      }'

# 用户登录
curl -X POST "http://localhost:8080/sessions" \
  -H "Content-Type: application/json" \
  -d '{
        "name": "张三",
        "tenant_id": "019a829b-2231-7740-ba47-243aca9eadeb",
        "password": "123456"
      }'

# 创建短链
curl -X POST "http://localhost:8080/shortlinks" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJNIjp7InRlbmFudF9pZCI6IjAxOWE4MjliLTIyMzEtNzc0MC1iYTQ3LTI0M2FjYTllYWRlYiIsInVzZXJfaWQiOiIwMTlhODI5Yy02YjNlLTc1NzYtOTY3ZC0zMDIwMmZlMGU1MmQifSwiaXNzIjoiVGVuYW50IFNlcnZpY2UiLCJleHAiOjE3NjM3MzI3MTR9.kdsez_ECfgu_gNjFDM2XeM0GuNA1_G9qPn3nc4q54JQ" \
  -d '{
        "original_url": "www.baidu.com",
        "expiration": 3600
      }'

# 短链解析和跳转
curl -X GET "http://localhost:8080/resolve/6tSb4nAJRB4"

# 查看今日排行榜
curl -X GET "localhost:8080/ranking/today?num=100" \
    -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJNIjp7InRlbmFudF9pZCI6IjAxOWE4MjliLTIyMzEtNzc0MC1iYTQ3LTI0M2FjYTllYWRlYiIsInVzZXJfaWQiOiIwMTlhODI5Yy02YjNlLTc1NzYtOTY3ZC0zMDIwMmZlMGU1MmQifSwiaXNzIjoiVGVuYW50IFNlcnZpY2UiLCJleHAiOjE3NjM3MzI3MTR9.kdsez_ECfgu_gNjFDM2XeM0GuNA1_G9qPn3nc4q54JQ"
```

```bash
输入ctrl+C关闭,或者输入`make prod-app-down`和`make prod-infra-down`关闭
```

# 架构介绍
客户端 -> api-gateway

api-gateway -> tenant service

api-gateway -> link service -> rabbimq,mysql

rabbitmq -> analytics service -> clickhouse

api-gateway -> analytics service

# TODO
## 功能
- [ ] 批量创建短链
- [ ] 增加一些确实的分析服务?

## 性能
- [ ] 更换消息队列为Kafka
- [ ] 将统计服务channel更换为spsc无锁队列提升性能
- [ ] 提升承载能力