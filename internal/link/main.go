package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/maypok86/otter/v2"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"github.com/turbo514/shortenurl-v2/link/adapter"
	"github.com/turbo514/shortenurl-v2/link/adapter/cache"
	"github.com/turbo514/shortenurl-v2/link/entity"
	"github.com/turbo514/shortenurl-v2/shared/migrate_helper"
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

	dbtx, err := sql.Open(dbtype, fmt.Sprintf(mysqldsn, cfg.DatabaseConfig.Username, cfg.DatabaseConfig.Password, cfg.DatabaseConfig.Host, cfg.DatabaseConfig.Port, cfg.DatabaseConfig.Dbname, cfg.DatabaseConfig.Options))
	if err != nil {
		slog.Error("连接数据库失败: %w", err)
		return
	}
	defer dbtx.Close()
	queries := model.New(dbtx)
	db := mysqldb.NewMysqlShortLinkDB(queries)

	//
	if err := migrate_helper.Up(
		cfg.DatabaseConfig.MigrateFilePath,
		fmt.Sprintf(migratedsn, cfg.DatabaseConfig.Username, cfg.DatabaseConfig.Password, cfg.DatabaseConfig.Host, cfg.DatabaseConfig.Port, cfg.DatabaseConfig.Dbname, cfg.DatabaseConfig.Options),
	); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			slog.Error("数据库迁移失败: %w", err)
		}
	}

	l1cache, err := cache.NewShortLinkL1Cache(&otter.Options[string, entity.ShortLink]{
		MaximumSize:      10000,                                                            // 缓存最多 1 万条记录
		ExpiryCalculator: otter.ExpiryAccessing[string, entity.ShortLink](5 * time.Minute), // 5 分钟不访问过期
	})
	if err != nil {
		slog.Error("L1Cache初始化失败", err, err.Error())
		return
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Common.Redis.Host, cfg.Common.Redis.Port),
	})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		slog.Error("redis连接创建失败", "err", err.Error())
		return
	}
	defer redisClient.Close()
	l2cache := cache.NewShortLinkL2Cache(redisClient)

	repo := adapter.NewShortLinkRepository(db, l1cache, l2cache)

	amqpconn, err := amqp091.Dial(fmt.Sprintf(rabbitmqurl, cfg.Common.Mq.Username, cfg.Common.Mq.Password, cfg.Common.Mq.Host, cfg.Common.Mq.Port))
	if err != nil {
		slog.Error("连接消息队列失败", "err", err.Error())
		return
	}
	defer amqpconn.Close()
	publisher := adapter.NewEventPublisher(amqpconn)
	if err := publisher.Init(); err != nil {
		slog.Error("初始化消息队列失败", "err", err.Error())
		return
	}

	service := usecase.NewLinkUseCase(repo, publisher)

	serviceServer := server.NewGrpcServer(service)
	grpcServer := grpc.NewServer()
	linkpb.RegisterLinkServiceServer(grpcServer, serviceServer)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.Port))
	if err != nil {
		slog.Error("监听端口失败", err, err.Error())
	}

	if err := grpcServer.Serve(lis); err != nil {
		slog.Error("开启服务器失败", err, err.Error())
	}
}
