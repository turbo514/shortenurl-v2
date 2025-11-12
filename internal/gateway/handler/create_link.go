package handler

import (
	"github.com/gin-gonic/gin"
	appcontext "github.com/turbo514/shortenurl-v2/gateway/app_context"
	"github.com/turbo514/shortenurl-v2/gateway/dto"
	"github.com/turbo514/shortenurl-v2/gateway/middleware"
	linkpb "github.com/turbo514/shortenurl-v2/shared/gen/proto/link"
	"github.com/turbo514/shortenurl-v2/shared/mylog"
	mytrace "github.com/turbo514/shortenurl-v2/shared/trace"
	"google.golang.org/protobuf/types/known/durationpb"
	"net/http"
)

func CreateLink(app *appcontext.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		tr := mytrace.GetTracer()
		ctx, span := tr.Start(c.Request.Context(), "ApiGateway.CreateLink")
		defer span.End()

		// 解析参数
		var req dto.CreateLinkRequest
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"errmsg":  err.Error(),
				"errcode": 100001,
			})
			return
		}
		mylog.GetLogger().Debug("获取到CreateLink请求", "req", req)
		tenantId, _ := c.Get(middleware.TENANTIDTYPE)
		userId, _ := c.Get(middleware.USERIDTYPE)

		// 向下游服务发送请求
		resp, err := app.Services.Link.CreateLink(ctx, &linkpb.CreateLinkRequest{
			TenantId:    tenantId.(string),
			UserId:      userId.(string),
			OriginalUrl: req.OriginalUrl,
			Expiration:  &durationpb.Duration{Seconds: req.Expiration},
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"errmsg":  err.Error(),
				"errcode": 100002,
			})
			return
		}

		// 返回响应
		c.JSON(http.StatusOK, dto.CreateLinkResponse{
			OriginalUrl: resp.OriginalUrl,
			ShortCode:   resp.ShortCode,
		})
	}
}
