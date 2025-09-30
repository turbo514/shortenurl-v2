package main

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	viper "github.com/turbo514/shortenurl-v2/shared/commonconfig"
	tenantpb "github.com/turbo514/shortenurl-v2/shared/gen/proto/tenant"
	"github.com/turbo514/shortenurl-v2/shared/migrate_helper"
	"github.com/turbo514/shortenurl-v2/tenant/config"
	"github.com/turbo514/shortenurl-v2/tenant/controller"
	"github.com/turbo514/shortenurl-v2/tenant/dao/repository"
	"github.com/turbo514/shortenurl-v2/tenant/service"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"log/slog"
	"net"
)

const mysqldsn = "%s:%s@tcp(%s:%d)/%s?%s"
const migratedsn = "mysql://%s:%s@tcp(%s:%d)/%s?%s"

func main() {
	v, err := viper.NewViper("global", "../shared/commonconfig/", "config", "./config/")
	if err != nil {
		panic(err)
	}
	cfg, err := config.NewConfig(v)
	if err != nil {
		panic(err)
	}

	//debug
	log.Printf("%+v\n", cfg)

	dsn := fmt.Sprintf(mysqldsn, cfg.Database.Username, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Dbname, cfg.Database.Options)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("连接数据库失败", "err", err.Error())
		return
	}

	// 数据库迁移
	if err := migrate_helper.Up(
		cfg.Database.MigrationFilePath,
		fmt.Sprintf(migratedsn, cfg.Database.Username, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Dbname, cfg.Database.Options),
	); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			slog.Error("数据库迁移失败", "err", err.Error())
			return
		}
	}

	userRepo := repository.NewUserRepo(db)
	tenantRepo := repository.NewTenantRepo(db)
	tenantService := service.NewTenantService(tenantRepo, userRepo)

	tokenService := service.NewTokenService("", "")

	handler := controller.NewHandler(tenantService, tokenService)

	server := grpc.NewServer()
	defer server.GracefulStop()
	tenantpb.RegisterTenantServiceServer(server, handler)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	if err := server.Serve(lis); err != nil {
		slog.Error("grpc服务器启动失败", "err", err.Error())
	}
}
