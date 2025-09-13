package appcontext

import "github.com/turbo514/shortenurl-v2/gateway/config"

type AppContext struct {
	Cfg      *config.Config
	Services *Services
}

func NewAppContext(cfg *config.Config, services *Services) *AppContext {
	return &AppContext{
		Cfg:      cfg,
		Services: services,
	}
}
