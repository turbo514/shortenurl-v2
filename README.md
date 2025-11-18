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

可以做到2700qps

由于是本机测试,应该有一定误差
一大部分性能损耗来自消息队列和序列化

几乎全部命中缓存

# 快速开始
```bash
git clone https://github.com/turbo514/shortenurl-v2.git

cd shortenurl-v2

make prod-infra-up

make kafka-init

make prod-app-up
```

```bash
# 创建租户
https POST "https://localhost:8080/tenants" name="张三的公司" --verify n

# 创建用户
https POST "https://localhost:8080/users" name="张三" tenant_id="019a91a2-57e9-735f-aa27-ec028ec57bb5" api_key="07e50f06-338e-  1f-8ee4-a5f38fad82e9" password="123456" --json --verify no

# 用户登录
https POST "https://localhost:8080/sessions" name="张三" tenant_id="019a91a2-57e9-735f-aa27-ec028ec57bb5" password="123456" --json --verify no

# 创建短链
https POST "https://localhost:8080/shortlinks" \
Authorization:"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJNIjp7InRlbmFudF9pZCI6IjAxOWE5MWEyLTU3ZTktNzM1Zi1hYTI3LWVjMDI4ZWM1N2JiNSIsInVzZXJfaWQiOiIwMTlhOTFhMi1hOTAwLTc2Y2UtODZiNy05ZDkxOTMyMjU5NWQifSwiaXNzIjoiVGVuYW50IFNlcnZpY2UiLCJleHAiOjE3NjQwNjIyOTV9.GBWW7Uqe-2zqEOuW-KoFaVpZRp7_hu5TrwNwa9h6Miw" \
original_url="www.baidu.com" expiration:=3600 \
--json --verify no

# 短链解析和跳转
https GET "https://localhost:8080/resolve/6tSb4nAJRB4"  --verify no

# 查看今日排行榜
https GET "localhost:8080/ranking/today?num=100" \
Authorization:"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJNIjp7InRlbmFudF9pZCI6IjAxOWE5MWEyLTU3ZTktNzM1Zi1hYTI3LWVjMDI4ZWM1N2JiNSIsInVzZXJfaWQiOiIwMTlhOTFhMi1hOTAwLTc2Y2UtODZiNy05ZDkxOTMyMjU5NWQifSwiaXNzIjoiVGVuYW50IFNlcnZpY2UiLCJleHAiOjE3NjQwNjIyOTV9.GBWW7Uqe-2zqEOuW-KoFaVpZRp7_hu5TrwNwa9h6Miw" \
--json --verify no
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