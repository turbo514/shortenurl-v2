package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oklog/run"
	"github.com/redis/go-redis/v9"
	appcontext "github.com/turbo514/shortenurl-v2/gateway/app_context"
	"github.com/turbo514/shortenurl-v2/gateway/config"
	"github.com/turbo514/shortenurl-v2/gateway/router"
	"github.com/turbo514/shortenurl-v2/shared/client"
	viper "github.com/turbo514/shortenurl-v2/shared/commonconfig"
	"github.com/turbo514/shortenurl-v2/shared/mylog"
	"github.com/turbo514/shortenurl-v2/shared/rate_limiter"
	mytrace "github.com/turbo514/shortenurl-v2/shared/trace"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"syscall"
)

func main() {
	// --- 日志设置 (slog) ---
	logger := mylog.GetLogger()

	ctx := context.Background()

	// 读取配置
	v, err := viper.NewViper("global", "../shared/commonconfig/", "config", "./config/")
	if err != nil {
		logger.Error("读取配置失败", "err", err.Error())
		return
	}
	cfg, err := config.NewConfig(v)
	if err != nil {
		logger.Error("初始化配置失败", "err", err.Error())
		return
	}
	logger.Debug("检查配置内容", "config", cfg)

	// --- 初始化 OpenTelemetry 追踪设置 ---
	if _, err := mytrace.InitOpenTelemetry(ctx, &cfg.Jaeger, &cfg.ServiceInfo); err != nil {
		logger.Error("初始化Trace Provider失败", "err", err.Error())
		return
	}

	// --- 初始化Redis连接 ---
	redisClient := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf(":%d", cfg.Redis.Port),
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
	})
	if err := redisClient.Ping(ctx).Err(); err != nil {
		logger.Error("创建Redis 连接失败", "err", err.Error())
		return
	}

	// --- 初始化限流器 ---
	globalRateLimiter := rate_limiter.NewRedisRateLimiter(redisClient, cfg.GlobalRateLimiter.Rate, cfg.GlobalRateLimiter.Capacity)
	localRateLimiter := rate.NewLimiter(rate.Limit(cfg.LocalRateLimiter.Rate), int(cfg.LocalRateLimiter.Capacity))

	// --- 初始化下游服务连接 ---
	services := &appcontext.Services{}
	if conn, err := client.NewLinkConn(
		fmt.Sprintf("%s:%d", cfg.LinkService.Host, cfg.LinkService.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	); err != nil {
		logger.Error("初始化与Link Service的连接失败", "err", err.Error())
		return
	} else {
		services.Link = client.NewLinkClient(conn)
	}
	if conn, err := client.NewTenantConn(
		fmt.Sprintf("%s:%d", cfg.TenantService.Host, cfg.TenantService.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	); err != nil {
		logger.Error("初始化与Tenant Service的连接失败")
		return
	} else {
		services.Tenant = client.NewTenantClient(conn)
	}
	if conn, err := client.NewAnalyticsConn(
		fmt.Sprintf("%s:%d", cfg.AnalyticsService.Host, cfg.AnalyticsService.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	); err != nil {
		logger.Error("初始化与Analytics Service的连接失败")
		return
	} else {
		services.Analytics = client.NewAnalyticsClient(conn)
	}
	app := appcontext.NewAppContext(cfg, services, globalRateLimiter, localRateLimiter)

	// --- 初始化路由逻辑 ---
	r := router.NewRouter(app)

	// 编排任务
	g := &run.Group{}

	// 1. 启动HTTP服务器,监听用户请求
	g.Add(func() error {
		gin.SetMode(gin.ReleaseMode)
		logger.Info("正在启动HTTP服务器", "port", cfg.Server.Port)
		if err := r.Run(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
			return fmt.Errorf("启动HTTP服务器失败", "err", err.Error())
		}
		return nil
	}, func(err error) {
	})

	// 2. 信号处理任务: 优雅关闭
	g.Add(run.SignalHandler(ctx, syscall.SIGINT, syscall.SIGTERM))

	// 运行所有任务
	if err := g.Run(); err != nil {
		logger.Error("程序退出", "err", err.Error())
		return
	}
}
