package handler

import (
	"github.com/gin-gonic/gin"
	appcontext "github.com/turbo514/shortenurl-v2/gateway/app_context"
	"github.com/turbo514/shortenurl-v2/gateway/dto"
	tenantpb "github.com/turbo514/shortenurl-v2/shared/gen/proto/tenant"
	"github.com/turbo514/shortenurl-v2/shared/mylog"
	mytrace "github.com/turbo514/shortenurl-v2/shared/trace"
	"net/http"
)

func Login(app *appcontext.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		tr := mytrace.GetTracer()
		ctx, span := tr.Start(c.Request.Context(), "ApiGateway.Login")
		defer span.End()

		logger := mylog.GetLogger()

		// 解析参数
		var req dto.LoginRequest
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"errmsg":  err.Error(),
				"errcode": 100001,
			})
			return
		}
		logger.Debug("接收到Login请求", "req", req)

		// 用户登录
		resp, err := app.Services.Tenant.Login(ctx, &tenantpb.LoginRequest{
			TenantId: req.TenantId,
			Name:     req.Name,
			Password: req.Password,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"errmsg":  err.Error(),
				"errcode": 100001,
			})
			return
		}

		// 响应用户请求,返回Token
		c.JSON(http.StatusOK, dto.LoginResponse{
			Token: resp.Token,
		})
	}
}
