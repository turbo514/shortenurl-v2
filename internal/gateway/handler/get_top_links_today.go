package handler

import (
	"github.com/gin-gonic/gin"
	appcontext "github.com/turbo514/shortenurl-v2/gateway/app_context"
	"github.com/turbo514/shortenurl-v2/gateway/dto"
	"github.com/turbo514/shortenurl-v2/gateway/middleware"
	analyticspb "github.com/turbo514/shortenurl-v2/shared/gen/proto/analytics"
	"github.com/turbo514/shortenurl-v2/shared/mylog"
	mytrace "github.com/turbo514/shortenurl-v2/shared/trace"
	"net/http"
)

func GetTopLinksToday(app *appcontext.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		tr := mytrace.GetTracer()
		ctx, span := tr.Start(c.Request.Context(), "ApiGateway.GetTopLinksToday")
		defer span.End()

		// 解析参数
		var req dto.GetTopLinksTodayRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"errmsg":  err.Error(),
				"errcode": 100001,
			})
			return
		}
		tenantId, _ := c.Get(middleware.TENANTIDTYPE)
		mylog.GetLogger().Debug("获取到GetTopLinksToday请求", "req", req, "tenant_id", tenantId)

		// 向下游服务发送请求
		resp, err := app.Services.Analytics.GetTopN(ctx, &analyticspb.GetTopNRequest{
			N:        req.Num,
			TenantId: tenantId.(string),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"errmsg":  err.Error(),
				"errcode": 100002,
			})
			return
		}

		// 返回响应
		topLinks := make([]dto.TopLinks, len(resp.TopLinks))
		for i := range resp.TopLinks {
			topLinks[i] = dto.TopLinks{
				ID:          resp.TopLinks[i].Id,
				OriginalUrl: resp.TopLinks[i].OriginalUrl,
				ClickTimes:  resp.TopLinks[i].ClickTimes,
			}
		}
		c.JSON(http.StatusOK, topLinks)
	}
}
