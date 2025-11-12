package middleware

import (
	"github.com/gin-gonic/gin"
	appcontext "github.com/turbo514/shortenurl-v2/gateway/app_context"
	"github.com/turbo514/shortenurl-v2/shared/keys"
	mytrace "github.com/turbo514/shortenurl-v2/shared/trace"
	"net/http"
)

func RateLimitMiddleware(app *appcontext.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := keys.GlobalRateKey
		if ok, err := app.GlobalRateLimiter.Allow(c.Request.Context(), key, 1); err == nil {
			if !ok {
				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
					"errmsg": "请求过于频繁,请稍后再试",
				})
				return
			} else {
				c.Next()
				return
			}
		} else {
			// 降级处理
			tr := mytrace.GetTracer()
			ctx, span := tr.Start(c.Request.Context(), "Middleware.RateLimit")
			defer span.End()

			span.RecordError(err)
			c.Request = c.Request.WithContext(ctx)
			if ok := app.LocalRateLimiter.Allow(); !ok {
				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
					"errmsg": "请求过于频繁,请稍后再试",
				})
				return
			} else {
				c.Next()
				return
			}
		}
	}
}
