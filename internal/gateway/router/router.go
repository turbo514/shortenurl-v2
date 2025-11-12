package router

import (
	"github.com/gin-gonic/gin"
	appcontext "github.com/turbo514/shortenurl-v2/gateway/app_context"
	"github.com/turbo514/shortenurl-v2/gateway/handler"
	"github.com/turbo514/shortenurl-v2/gateway/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func NewRouter(app *appcontext.AppContext) *gin.Engine {
	r := gin.Default()

	r.Use(otelgin.Middleware("api-gateway"), middleware.RateLimitMiddleware(app))
	r.POST("/sessions", handler.Login(app))
	r.POST("/tenants", handler.CreateTenant(app))
	r.POST("/users", handler.CreateUser(app))
	r.GET("/resolve/:short_code", handler.ResolveLink(app))

	auth := r.Group("/", middleware.AuthMiddleware(app))
	auth.POST("/shortlinks", handler.CreateLink(app))
	auth.GET("/ranking/today", handler.GetTopLinksToday(app))

	return r
}
