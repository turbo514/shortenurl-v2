package handler

import (
	"github.com/gin-gonic/gin"
	appcontext "github.com/turbo514/shortenurl-v2/gateway/app_context"
	"github.com/turbo514/shortenurl-v2/gateway/dto"
	tenantpb "github.com/turbo514/shortenurl-v2/shared/gen/proto/tenant"
	"net/http"
)

func CreateTenant(app *appcontext.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CreateTenantRequest
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			// TODO: 以后再统一格式
			c.JSON(http.StatusBadRequest, gin.H{
				"errmsg":  err.Error(),
				"errcode": 100001,
			})
			return
		}

		// TODO: 参数检验

		//context.WithTimeout(c.Request.Context(),)
		resp, err := app.Services.Tenant.CreateTenant(c.Request.Context(), &tenantpb.CreateTenantRequest{
			Name: req.Name,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"errmsg":  err.Error(),
				"errcode": 100002,
			})
		} else {
			c.JSON(http.StatusOK, dto.CreateTenantResponse{
				TenantId: resp.TenantId,
				ApiKey:   resp.ApiKey,
			})
		}
	}
}
