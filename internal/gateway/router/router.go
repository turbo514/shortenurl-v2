package router

import (
	"github.com/gin-gonic/gin"
	appcontext "github.com/turbo514/shortenurl-v2/gateway/app_context"
	"github.com/turbo514/shortenurl-v2/gateway/handler"
)

func NewRouter(ctx *appcontext.AppContext) *gin.Engine {
	r := gin.Default()

	r.POST("/tenants", handler.CreateTenant(ctx))
	r.POST("/users", handler.CreateUser(ctx))
	r.POST("/sessions", handler.Login(ctx))

	return r
}
