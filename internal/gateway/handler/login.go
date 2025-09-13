package handler

import (
	"github.com/gin-gonic/gin"
	appcontext "github.com/turbo514/shortenurl-v2/gateway/app_context"
	"github.com/turbo514/shortenurl-v2/gateway/dto"
	tenantpb "github.com/turbo514/shortenurl-v2/shared/gen/proto/tenant"
	"net/http"
)

func Login(app *appcontext.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.LoginRequest
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"errmsg":  err.Error(),
				"errcode": 100001,
			})
			return
		}

		ctx := c.Request.Context()
		resp, err := app.Services.Tenant.Login(ctx, &tenantpb.LoginRequest{
			Name:     req.Name,
			Password: req.Password,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"errmsg":  err.Error(),
				"errcode": 100001,
			})
		} else {
			c.JSON(http.StatusOK, dto.LoginResponse{
				Token: resp.Token,
			})
		}
	}
}
