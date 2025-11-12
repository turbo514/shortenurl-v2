package appcontext

import (
	"github.com/turbo514/shortenurl-v2/gateway/config"
	"github.com/turbo514/shortenurl-v2/shared/rate_limiter"
	"golang.org/x/time/rate"
)

type AppContext struct {
	Cfg               *config.Config
	Services          *Services
	GlobalRateLimiter rate_limiter.IRateLimiter
	LocalRateLimiter  *rate.Limiter
}

func NewAppContext(cfg *config.Config, services *Services, globalRateLimiter rate_limiter.IRateLimiter, localRateLimiter *rate.Limiter) *AppContext {
	return &AppContext{
		Cfg:               cfg,
		Services:          services,
		GlobalRateLimiter: globalRateLimiter,
		LocalRateLimiter:  localRateLimiter,
	}
}
