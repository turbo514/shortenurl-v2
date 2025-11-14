package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus" // Prometheus 指标提供者
	"github.com/maypok86/otter/v2"
	"github.com/oklog/run" // 用于管理多个 Goroutine 任务（如 gRPC server, HTTP server）
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/turbo514/shortenurl-v2/link/config"
	"github.com/turbo514/shortenurl-v2/link/domain"
	otterrepo "github.com/turbo514/shortenurl-v2/link/infrastructure/otter_repository"
	"github.com/turbo514/shortenurl-v2/link/infrastructure/rabbitmq_publisher"
	redisrepo "github.com/turbo514/shortenurl-v2/link/infrastructure/redis_repository"
	"github.com/turbo514/shortenurl-v2/link/metrics"
	"github.com/turbo514/shortenurl-v2/shared/mylog"
	myprom "github.com/turbo514/shortenurl-v2/shared/prom"
	mytrace "github.com/turbo514/shortenurl-v2/shared/trace"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"net/http"
	"syscall"
	"time"

	"github.com/turbo514/shortenurl-v2/link/infrastructure/grpc_server"
	mysqlrepo "github.com/turbo514/shortenurl-v2/link/infrastructure/mysql_repository"
	"github.com/turbo514/shortenurl-v2/link/infrastructure/mysql_repository/model"
	"github.com/turbo514/shortenurl-v2/link/usecase"
	"github.com/turbo514/shortenurl-v2/shared/commonconfig"
	linkpb "github.com/turbo514/shortenurl-v2/shared/gen/proto/link"
	"google.golang.org/grpc"

	"net"
)

func main() {
	// --- 日志设置 (slog) ---
	logger := mylog.GetLogger()
	// 针对 RPC 请求创建带默认字段的日志器
	//rpcLogger := logger.With("service", "gRPC/server", "component", component)

	ctx := context.Background()

	// --- 获取配置 ---
	v, err := commonconfig.NewViper(commonconfig.GlobalFile, commonconfig.GlobalPath, commonconfig.ServiceFile, commonconfig.ServicePath)
	if err != nil {
		logger.Error("读取配置失败", "err", err.Error())
		return
	}
	cfg, err := config.NewConfig(v)
	if err != nil {
		logger.Error("初始化配置失败", "err", err.Error())
		return
	}
	logger.Debug("测试配置文件载入", "配置文件", cfg)

	// --- OpenTelemetry 追踪设置 ---
	if _, err := mytrace.InitOpenTelemetry(
		ctx,
		&cfg.Jaeger,
		&cfg.ServiceInfo,
	); err != nil {
		logger.Error("初始化Trace Provider失败", "err", err.Error())
		return
	}

	// --- Prometheus 指标设置 ---
	reg := prometheus.NewRegistry()
	srvMetrics := grpcprom.NewServerMetrics(
		grpcprom.WithServerHandlingTimeHistogram(
			grpcprom.WithHistogramOpts(&prometheus.HistogramOpts{
				Name:    "link_service_request_duration_seconds",
				Help:    "短链接服务的grpc请求响应用时直方图",
				Buckets: prometheus.DefBuckets, // 小于1毫秒,小于10毫秒,小于100毫秒,小于1秒,小于10秒,,更大的
			}),
		),
	)
	reg.MustRegister(srvMetrics)
	reg.MustRegister(metrics.GetListMetrics())

	// --- MySQL连接初始化 ---
	dbtx, err := sql.Open(
		"mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", cfg.Mysql.Username, cfg.Mysql.Password, cfg.Mysql.Host, cfg.Mysql.Port, cfg.Mysql.DbName, cfg.Mysql.Options),
	)
	if err != nil {
		logger.Error("连接数据库失败: %w", err)
		return
	}
	defer dbtx.Close()

	// --- 数据库迁移 ---
	//if err := migrate_helper.Up(
	//	"./infrastructure/migrations/",
	//	fmt.Sprintf("mysql://%s:%s@tcp(%s:%d)/%s?%s", cfg.Mysql.Username, cfg.Mysql.Password, cfg.Mysql.Host, cfg.Mysql.Port, cfg.Mysql.DbName, cfg.Mysql.Options),
	//); err != nil {
	//	if !errors.Is(err, migrate.ErrNoChange) {
	//		logger.Error("数据库迁移失败: %w", err)
	//	}
	//}

	// --- 初始化本地缓存 ---
	otterCache, err := otter.New(&otter.Options[string, *domain.ShortLink]{
		MaximumSize:      1000,                                                        // 缓存最多 1 千条记录
		ExpiryCalculator: otter.ExpiryWriting[string, *domain.ShortLink](time.Minute), // 1分钟不访问过期
	})
	if err != nil {
		logger.Error("otter cache创建失败", "err", err.Error())
		return
	}

	// --- 初始化Redis连接 ---
	redisClient := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
	})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		logger.Error("redis连接创建失败", "err", err.Error())
		return
	}
	defer redisClient.Close()
	// 开启 tracing instrumentation
	if err := redisotel.InstrumentTracing(redisClient); err != nil {
		logger.Error("redis开启tracing instrumentation失败", "err", err.Error())
		return
	}
	// 开启 metrics instrumentation
	if err := redisotel.InstrumentMetrics(redisClient); err != nil {
		logger.Error("redis开启metrics instrumentation失败", "err", err.Error())
		return
	}

	// --- 初始化repository层 ---
	queries := model.New(dbtx)
	mysqlRepo := mysqlrepo.NewMysqlShortLinkDB(queries)
	redisRepo := redisrepo.NewRedisCacheRepository(redisClient, mysqlRepo)
	otterRepo := otterrepo.NewOtterCacheRepository(otterCache, redisRepo)

	// --- 初始化RabbitMQ消息队列 ---
	amqpconn, err := amqp091.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d", cfg.RabbitMq.Username, cfg.RabbitMq.Password, cfg.RabbitMq.Host, cfg.RabbitMq.Port))
	if err != nil {
		logger.Error("连接消息队列失败", "err", err.Error())
		return
	}
	defer amqpconn.Close()

	// --- 初始化事件发送器 ---
	publisher := rabbitmq_publisher.NewEventPublisher(amqpconn)
	if err := publisher.Init(); err != nil {
		logger.Error("初始化消息队列失败", "err", err.Error())
		return
	}

	// --- 初始化服务逻辑 ---
	service := usecase.NewLinkUseCase(otterRepo, publisher)

	// --- 创建grpc服务器 ---
	gServer := grpc.NewServer(
		// 启用otel追踪
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		// 监控
		grpc.ChainUnaryInterceptor(
			srvMetrics.UnaryServerInterceptor(
				grpcprom.WithExemplarFromContext(myprom.ExemplarFromContext),
				grpcprom.WithLabelsFromContext(myprom.LabelsFromContext),
			),
			metrics.UnaryServerInterceptor(),
		),
	)

	// --- 注册grpc服务实现 ---
	server := grpc_server.NewGrpcServer(service)
	linkpb.RegisterLinkServiceServer(gServer, server)

	// --- 初始化 Prometheus 指标，为所有已注册的方法预先创建指标 ---
	srvMetrics.InitializeMetrics(gServer)

	// --- 启动/停止 Goroutine 管理 ---
	g := &run.Group{}

	// 1. gRPC 服务器启动/停止任务
	g.Add(func() error {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.Port))
		if err != nil {
			logger.Error("监听grpc端口失败", "err", err.Error())
			return err
		}
		logger.Info("正在启动grpc服务器", "addr", lis.Addr().String())
		if err := gServer.Serve(lis); err != nil {
			logger.Error("开启grpc服务器失败", "err", err.Error())
			return err
		}
		return nil
	}, func(err error) {
		gServer.GracefulStop()
		gServer.Stop()
	})

	// 2. HTTP 指标服务器启动/停止任务
	httpSrv := &http.Server{Addr: fmt.Sprintf(":%d", cfg.Prometheus.Port)}
	g.Add(func() error {
		m := http.NewServeMux()
		m.Handle("/metrics", promhttp.HandlerFor(
			reg, // 使用自定义的注册表
			promhttp.HandlerOpts{
				EnableOpenMetrics: true,
			},
		))
		httpSrv.Handler = m
		logger.Info("正在启动http服务器", "addr", httpSrv.Addr)
		if err := httpSrv.ListenAndServe(); err != nil {
			logger.Error("启动http服务器失败", "err", err.Error())
			return err
		}
		return nil
	}, func(err error) {
		if err := httpSrv.Close(); err != nil {
			logger.Error("关闭http服务器失败", "err", err.Error())
		}
	})

	// 3. 信号处理任务：优雅关闭
	g.Add(run.SignalHandler(context.Background(), syscall.SIGINT, syscall.SIGTERM))

	// 运行所有任务，直到其中一个失败或收到信号
	if err := g.Run(); err != nil {
		logger.Error("程序退出", "err", err.Error())
		return
	}
}
