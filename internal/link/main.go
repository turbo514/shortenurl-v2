package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus" // Prometheus 指标提供者
	"github.com/maypok86/otter/v2"
	"github.com/oklog/run" // 用于管理多个 Goroutine 任务（如 gRPC server, HTTP server）
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"github.com/turbo514/shortenurl-v2/link/adapter"
	"github.com/turbo514/shortenurl-v2/link/adapter/cache"
	"github.com/turbo514/shortenurl-v2/link/entity"
	"github.com/turbo514/shortenurl-v2/shared/migrate_helper"
	myprom "github.com/turbo514/shortenurl-v2/shared/prom"
	mytrace "github.com/turbo514/shortenurl-v2/shared/trace"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"net/http"
	"syscall"
	"time"

	"github.com/turbo514/shortenurl-v2/link/infrastructure/config"
	"github.com/turbo514/shortenurl-v2/link/infrastructure/mysqldb"
	"github.com/turbo514/shortenurl-v2/link/infrastructure/mysqldb/model"
	"github.com/turbo514/shortenurl-v2/link/infrastructure/server"
	"github.com/turbo514/shortenurl-v2/link/usecase"
	"github.com/turbo514/shortenurl-v2/shared/commonconfig"
	linkpb "github.com/turbo514/shortenurl-v2/shared/gen/proto/link"
	"google.golang.org/grpc"

	"log/slog"
	"net"
)

const dbtype = "mysql"
const mysqldsn = "%s:%s@tcp(%s:%d)/%s?%s"
const rabbitmqurl = "amqp://%s:%s@%s:%d"
const migratedsn = "mysql://%s:%s@tcp(%s:%d)/%s?%s"

func main() {
	ctx := context.Background()

	// 获取配置
	v, err := commonconfig.NewViper("global", "../shared/commonconfig/", "config", "./infrastructure/config/")
	if err != nil {
		slog.Error("读取配置失败", "err", err.Error())
		return
	}
	cfg, err := config.NewConfig(v)
	if err != nil {
		slog.Error("初始化配置失败", "err", err.Error())
		return
	}

	fmt.Printf("%+v\n", cfg)

	// 初始化otel追踪
	if _, err := mytrace.InitOpenTelemetry(
		&cfg.Common.Jaeger,
		&cfg.ServiceInfo,
		ctx,
	); err != nil {
		slog.Error("初始化Trace Provider失败", "err", err.Error())
		return
	}

	// 初始化prometheus监控
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

	// 初始化数据库连接
	dbtx, err := sql.Open(
		dbtype,
		fmt.Sprintf(mysqldsn, cfg.DatabaseConfig.Username, cfg.DatabaseConfig.Password, cfg.DatabaseConfig.Host, cfg.DatabaseConfig.Port, cfg.DatabaseConfig.Dbname, cfg.DatabaseConfig.Options),
	)
	if err != nil {
		slog.Error("连接数据库失败: %w", err)
		return
	}
	defer dbtx.Close()
	queries := model.New(dbtx)
	db := mysqldb.NewMysqlShortLinkDB(queries)

	// 数据库迁移
	if err := migrate_helper.Up(
		cfg.DatabaseConfig.MigrateFilePath,
		fmt.Sprintf(migratedsn, cfg.DatabaseConfig.Username, cfg.DatabaseConfig.Password, cfg.DatabaseConfig.Host, cfg.DatabaseConfig.Port, cfg.DatabaseConfig.Dbname, cfg.DatabaseConfig.Options),
	); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			slog.Error("数据库迁移失败: %w", err)
		}
	}

	// 初始化l1缓存
	l1cache, err := cache.NewShortLinkL1Cache(&otter.Options[string, entity.ShortLink]{
		MaximumSize:      10000,                                                            // 缓存最多 1 万条记录
		ExpiryCalculator: otter.ExpiryAccessing[string, entity.ShortLink](5 * time.Minute), // 5 分钟不访问过期
	})
	if err != nil {
		slog.Error("L1Cache初始化失败", err, err.Error())
		return
	}

	// 初始化分布式缓存
	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Common.Redis.Host, cfg.Common.Redis.Port),
	})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		slog.Error("redis连接创建失败", "err", err.Error())
		return
	}
	defer redisClient.Close()
	l2cache := cache.NewShortLinkL2Cache(redisClient)

	// 初始化repository
	repo := adapter.NewShortLinkRepository(db, l1cache, l2cache)

	// 初始化消息队列
	amqpconn, err := amqp091.Dial(fmt.Sprintf(rabbitmqurl, cfg.Common.Mq.Username, cfg.Common.Mq.Password, cfg.Common.Mq.Host, cfg.Common.Mq.Port))
	if err != nil {
		slog.Error("连接消息队列失败", "err", err.Error())
		return
	}
	defer amqpconn.Close()
	// 初始化事件发送器
	publisher := adapter.NewEventPublisher(amqpconn)
	if err := publisher.Init(); err != nil {
		slog.Error("初始化消息队列失败", "err", err.Error())
		return
	}

	// 初始化服务逻辑
	service := usecase.NewLinkUseCase(repo, publisher)

	// 注册grpc服务器
	serviceServer := server.NewGrpcServer(service)
	grpcServer := grpc.NewServer(
		// 启用Otel追踪
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		// 监控
		grpc.ChainUnaryInterceptor(
			srvMetrics.UnaryServerInterceptor(
				grpcprom.WithExemplarFromContext(myprom.ExemplarFromContext),
				grpcprom.WithLabelsFromContext(myprom.LabelsFromContext),
			),
		),
	)
	linkpb.RegisterLinkServiceServer(grpcServer, serviceServer)

	// 初始化Prometheus指标
	srvMetrics.InitializeMetrics(grpcServer)

	g := &run.Group{}
	g.Add(func() error {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.GrpcPort))
		if err != nil {
			return fmt.Errorf("监听grpc端口失败: %w", err)
		}
		slog.Info("正在启动grpc服务器", "addr", lis.Addr().String())
		if err := grpcServer.Serve(lis); err != nil {
			return fmt.Errorf("开启grpc服务器失败: %w", err)
		}
		return nil
	}, func(err error) {
		if err != nil {
			slog.Error("grpc服务器出错", "err", err.Error())
		}
		grpcServer.GracefulStop()
		grpcServer.Stop()
	})

	httpSrv := &http.Server{Addr: fmt.Sprintf(":%d", cfg.Server.HttpPort)}
	g.Add(func() error {
		m := http.NewServeMux()
		m.Handle("/metrics", promhttp.HandlerFor(
			reg,
			promhttp.HandlerOpts{
				EnableOpenMetrics: true,
			},
		))
		httpSrv.Handler = m
		slog.Info("正在启动http服务器", "addr", httpSrv.Addr)
		if err := httpSrv.ListenAndServe(); err != nil {
			return fmt.Errorf("启动http服务器失败: %w", err)
		}
		return nil
	}, func(err error) {
		if err != nil {
			slog.Error("http服务器出错", "err", err.Error())
		}
		if err := httpSrv.Close(); err != nil {
			slog.Error("关闭http服务器失败", "err", err.Error())
		}
	})

	// 信号处理任务：优雅关闭
	g.Add(run.SignalHandler(context.Background(), syscall.SIGINT, syscall.SIGTERM))

	// 运行所有任务，直到其中一个失败或收到信号
	if err := g.Run(); err != nil {
		slog.Error("程序退出", "err", err.Error())
		return
	}
}
