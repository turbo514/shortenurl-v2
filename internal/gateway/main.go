package main

import (
	"fmt"
	appcontext "github.com/turbo514/shortenurl-v2/gateway/app_context"
	"github.com/turbo514/shortenurl-v2/gateway/config"
	"github.com/turbo514/shortenurl-v2/gateway/router"
	"github.com/turbo514/shortenurl-v2/shared/client"
	viper "github.com/turbo514/shortenurl-v2/shared/commonconfig"
)

func main() {
	v, _ := viper.NewViper("global.yaml", "../shared/config/", "config.yaml", "./config/")
	cfg, _ := config.NewConfig(v)

	// debug
	fmt.Printf("%+v\n", cfg)

	services := &appcontext.Services{}

	if conn, err := client.NewLinkConn(fmt.Sprintf("%s:%d", cfg.Services.Link, cfg.Services.LinkPort)); err == nil {
		services.Link = client.NewLinkClient(conn)
	}

	if conn, err := client.NewTenantConn(fmt.Sprintf("%s:%d", cfg.Services.Link, cfg.Services.LinkPort)); err == nil {
		services.Tenant = client.NewTenantClient(conn)
	}

	if conn, err := client.NewAnalyticsConn(fmt.Sprintf("%s:%d", cfg.Services.Link, cfg.Services.LinkPort)); err == nil {
		services.Analytics = client.NewAnalyticsClient(conn)
	}

	app := appcontext.NewAppContext(cfg, services)

	r := router.NewRouter(app)

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
