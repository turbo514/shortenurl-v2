package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"github.com/turbo514/shortenurl-v2/analytics/config"
	"github.com/turbo514/shortenurl-v2/analytics/controller"
	"github.com/turbo514/shortenurl-v2/analytics/cqrs/command"
	"github.com/turbo514/shortenurl-v2/analytics/cqrs/query"
	"github.com/turbo514/shortenurl-v2/analytics/infra"
	"github.com/turbo514/shortenurl-v2/analytics/service"
	"github.com/turbo514/shortenurl-v2/analytics/util"
	"github.com/turbo514/shortenurl-v2/shared/commonconfig"
	analyticspb "github.com/turbo514/shortenurl-v2/shared/gen/proto/analytics"
	"github.com/turbo514/shortenurl-v2/shared/mylog"
	myprom "github.com/turbo514/shortenurl-v2/shared/prom"
	"github.com/turbo514/shortenurl-v2/shared/rabbitmq"
	mytrace "github.com/turbo514/shortenurl-v2/shared/trace"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
	"net"
	"net/http"
	"syscall"
)
import grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus" // Prometheus 指标提供者

const component = "grpc-example"

func main() {
	// --- 日志设置 (slog) ---
	logger := mylog.GetLogger()
	// 针对 RPC 请求创建带默认字段的日志器
	//rpcLogger := logger.With("service", "gRPC/server", "component", component)

	ctx := context.Background()

	// --- 读取配置 ---
	v, err := commonconfig.NewViper("global", "../shared/commonconfig/", "config", "./config/")
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

	// --- Prometheus 指标设置 ---
	srvMetrics := grpcprom.NewServerMetrics(
		// 启用服务器请求处理时间直方图 (Histogram)
		grpcprom.WithServerHandlingTimeHistogram(
			grpcprom.WithHistogramOpts(&prometheus.HistogramOpts{
				Name:    "analytics_service_request_duration_seconds",
				Help:    "分析服务的grpc请求相应用时直方图",
				Buckets: prometheus.DefBuckets,
			})),
	)
	reg := prometheus.NewPedanticRegistry() // 创建一个自定义的 Prometheus 注册表 (Registry)
	reg.MustRegister(srvMetrics)

	// --- OpenTelemetry 追踪设置 ---
	if _, err := mytrace.InitOpenTelemetry(ctx, &cfg.Jaeger, &cfg.ServiceInfo); err != nil {
		logger.Error("初始化Trace Provider失败", "err", err.Error())
		return
	}

	// --- ClickHouse连接 初始化 ---
	clickhouseConn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%d", cfg.Clickhouse.Host, cfg.Clickhouse.Port)},
		Auth: clickhouse.Auth{
			Database: cfg.Clickhouse.DbName,
			Username: cfg.Clickhouse.Username,
			Password: cfg.Clickhouse.Password,
		},
		Debug: true,
		Debugf: func(format string, v ...any) {
			logger.Info(fmt.Sprintf(format, v...))
		},
		TLS: nil,
	})
	if err != nil {
		logger.Error("连接ClickHouse失败", "err", err.Error())
		return
	}
	if err := clickhouseConn.Ping(ctx); err != nil {
		var exception *clickhouse.Exception
		if errors.As(err, &exception) {
			logger.Error("Ping ClickHouse失败", "code", exception.Code, "message", exception.Message, "stack", exception.StackTrace)
		}
	}

	// --- Redis连接 初始化 ---
	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
	})
	if err := redisClient.Ping(ctx).Err(); err != nil {
		logger.Error("redis连接创建失败", "err", err.Error())
		return
	}
	defer redisClient.Close()

	// --- MySQL连接 初始化 ---
	mysqlDb, err := gorm.Open(
		mysql.Open(
			fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", cfg.Mysql.Username, cfg.Mysql.Password, cfg.Mysql.Host, cfg.Mysql.Port, cfg.Mysql.DbName, cfg.Mysql.Options),
		))
	if err != nil {
		logger.Error("mysql连接创建失败", "err", err.Error())
		return
	}
	// --- 给数据库添加追踪 ---
	if err := mysqlDb.Use(tracing.NewPlugin(tracing.WithoutMetrics())); err != nil {
		logger.Error("gorm添加opentelemetry插件失败", "err", err.Error())
		return
	}

	// --- 消息队列 RabbitMq 初始化 ---
	amqpConn, err := amqp091.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d", cfg.RabbitMq.Username, cfg.RabbitMq.Password, cfg.RabbitMq.Host, cfg.RabbitMq.Port))
	if err != nil {
		logger.Error("链接消息队列失败", "err", err.Error())
		return
	}
	defer amqpConn.Close()
	// --- 队列初始化 ---
	if err := util.InitQueue(amqpConn); err != nil {
		logger.Error("初始化消息队列失败", "err", err.Error())
		return
	}

	// --- Database层初始化 ---
	clickhouseDb := infra.NewClickhouseDb(clickhouseConn)

	// --- 事件服务逻辑编排初始化 ---
	clickCounter := infra.NewRedisClickCounter(redisClient)

	createClickEventHandler := command.NewCreateClickEventHandler(clickhouseDb, clickCounter)
	getLinksHandler := query.NewGetLinksHandler(mysqlDb, redisClient)
	analyticsService := service.NewAnalyticsService(getLinksHandler, clickCounter)

	// --- 启动/停止 Goroutine 管理 (使用 oklog/run) ---
	g := &run.Group{}

	// 1. 监听链接点击事件
	clickEventReceiver := infra.NewRabbitMqClickEventReceiver(amqpConn, rabbitmq.ClickEventQueue)
	clickEventCh := clickEventReceiver.GetChannel()
	g.Add(func() error {
		logger.Info("正在监听点击事件")
		if err := clickEventReceiver.Start(); err != nil {
			logger.Error("监听链接点击事件失败", "err", err.Error())
			return err
		}
		return nil
	}, func(err error) {
		clickEventReceiver.Close()
		logger.Info("已停止监听点击事件")
	})

	// 2. 处理链接点击事件
	clickEventHandler := controller.NewClickEventHandler(clickEventCh, 2048, createClickEventHandler)
	g.Add(func() error {
		logger.Info("正在处理点解事件")
		if err := clickEventHandler.Start(); err != nil {
			logger.Error("处理链接点击事件失败", "err", err.Error())
		}
		return nil
	}, func(err error) {
		clickEventHandler.Close()
		logger.Info("已停止处理点击事件")
	})

	// 3. prometheus拉取用的http指标服务器
	httpSrv := http.Server{
		Addr: fmt.Sprintf(":%d", cfg.Prometheus.Port),
	}
	g.Add(func() error {
		m := http.NewServeMux()
		m.Handle("/metrics", promhttp.HandlerFor(
			reg,
			promhttp.HandlerOpts{
				EnableOpenMetrics: true,
			},
		))
		httpSrv.Handler = m
		logger.Info("正在启动Http服务器", "addr", httpSrv.Addr)
		if err := httpSrv.ListenAndServe(); err != nil {
			logger.Info("开启Http服务器失败", "err", err.Error())
			return err
		}
		return nil
	}, func(err error) {
		if err := httpSrv.Close(); err != nil {
			logger.Error("关闭http服务器失败", "err", err.Error())
		} else {
			logger.Info("关闭http服务器成功")
		}
	})

	// 4. grpc服务器,处理上游请求
	gServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			srvMetrics.UnaryServerInterceptor(
				grpcprom.WithExemplarFromContext(myprom.ExemplarFromContext),
				grpcprom.WithLabelsFromContext(myprom.LabelsFromContext),
			),
		),
	)
	server := controller.NewServiceHandler(analyticsService)
	analyticspb.RegisterAnalyticsServiceServer(gServer, server)
	g.Add(func() error {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.Port))
		if err != nil {
			logger.Error("监听Grpc端口失败", "err", err.Error())
			return err
		}
		defer lis.Close()

		logger.Info("正在启动grpc服务器", "port", lis.Addr())
		if err := gServer.Serve(lis); err != nil {
			logger.Error("启动Grpc服务器失败", "err", err.Error())
			return err
		}
		return nil
	}, func(err error) {
		gServer.GracefulStop()
		gServer.Stop()
		logger.Info("grpc服务器已关闭")
	})

	// 5. 信号处理任务：优雅关闭
	g.Add(run.SignalHandler(ctx, syscall.SIGINT, syscall.SIGTERM))

	// 运行所有服务
	if err := g.Run(); err != nil {
		logger.Error("程序报错", "err", err.Error())
		return
	}
}
