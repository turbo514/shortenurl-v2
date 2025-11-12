package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	appcontext "github.com/turbo514/shortenurl-v2/gateway/app_context"
	"github.com/turbo514/shortenurl-v2/gateway/dto"
	tenantpb "github.com/turbo514/shortenurl-v2/shared/gen/proto/tenant"
	"github.com/turbo514/shortenurl-v2/shared/mylog"
	mytrace "github.com/turbo514/shortenurl-v2/shared/trace"
	"github.com/turbo514/shortenurl-v2/shared/zerr"
	"net/http"
)

func CreateUser(app *appcontext.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		tr := mytrace.GetTracer()
		ctx, span := tr.Start(c.Request.Context(), "ApiGateway.CreateUser")
		defer span.End()

		// 解析参数
		// TODO: 参数校验
		var req dto.CreateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"errmsg":  err.Error(),
				"errcode": 100001,
			})
			return
		}
		mylog.GetLogger().Debug("接收到CreateUser请求", "req", req)

		// 向下游服务发送请求
		_, err := app.Services.Tenant.CreateUser(ctx, &tenantpb.CreateUserRequest{
			Name:     req.Name,
			Password: req.Password,
			TenantId: req.TenantId,
			ApiKey:   req.Apikey,
		})
		if err != nil {
			if errors.Is(err, zerr.ErrDuplicateEntry) {
				c.JSON(http.StatusConflict, gin.H{
					"errmsg": "该用户已存在",
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"errmsg":  err.Error(), // TODO: 内部错误
					"errcode": 100002,
				})
			}
			return
		}

		// 响应用户请求
		c.JSON(http.StatusOK, gin.H{})
	}
}
