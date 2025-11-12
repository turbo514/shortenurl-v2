package handler

import (
	"github.com/gin-gonic/gin"
	appcontext "github.com/turbo514/shortenurl-v2/gateway/app_context"
	"github.com/turbo514/shortenurl-v2/gateway/dto"
	linkpb "github.com/turbo514/shortenurl-v2/shared/gen/proto/link"
	"github.com/turbo514/shortenurl-v2/shared/mylog"
	mytrace "github.com/turbo514/shortenurl-v2/shared/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
	"net/http"
)

func ResolveLink(app *appcontext.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := mylog.GetLogger()

		tr := mytrace.GetTracer()
		ctx, span := tr.Start(c.Request.Context(), "ApiGateway.ResolveLink")
		defer span.End()

		// 解析参数
		var req dto.ResolveLinkRequest
		if err := c.ShouldBindUri(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"errmsg":  err.Error(),
				"errcode": 100001,
			})
			return
		}
		logger.Debug("接收到ResolveLink请求", "req", req)

		req2 := &linkpb.ResolveLinkRequest{
			ShortCode: req.ShortCode,
			IpAddress: net.ParseIP(c.ClientIP()),
			UserAgent: c.Request.UserAgent(),
			Referrer:  c.Request.Referer(),
		}
		logger.Debug("发送给链接解析服务的请求", "req", req2)
		resp, err := app.Services.Link.ResolveLink(ctx, req2)
		if err != nil {
			if st, ok := status.FromError(err); ok {
				if st.Code() == codes.NotFound {
					c.JSON(http.StatusNotFound, gin.H{
						"errmsg":  st.Message(),
						"errcode": 100002,
					})
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{
						"errmsg":  st.Message(),
						"errcode": 100003,
					})
				}
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"errmsg":  err.Error(), // TODO: 内部错误
					"errcode": 100004,
				})
			}
			return
		}

		// 响应用户请求,进行跳转
		c.Redirect(http.StatusTemporaryRedirect, "https://"+resp.OriginalUrl)
	}
}
