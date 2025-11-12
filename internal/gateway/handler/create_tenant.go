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

func CreateTenant(app *appcontext.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		tr := mytrace.GetTracer()
		ctx, span := tr.Start(c.Request.Context(), "ApiGateway.CreateTenant")
		defer span.End()

		// 解析参数
		// TODO: 参数检验
		var req dto.CreateTenantRequest
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			// TODO: 以后再统一格式
			c.JSON(http.StatusBadRequest, gin.H{
				"errmsg":  err.Error(),
				"errcode": 100001,
			})
			return
		}
		mylog.GetLogger().Debug("接收到CreateTenant请求", "req", req)

		// 向下游服务发送请求
		resp, err := app.Services.Tenant.CreateTenant(ctx, &tenantpb.CreateTenantRequest{
			Name: req.Name,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"errmsg":  err.Error(),
				"errcode": 100002,
			})
			return
		}

		// 返回用户响应
		c.JSON(http.StatusOK, dto.CreateTenantResponse{
			TenantId: resp.TenantId,
			ApiKey:   resp.ApiKey,
		})
	}
}
