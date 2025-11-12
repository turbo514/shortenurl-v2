package middleware

import (
	"github.com/gin-gonic/gin"
	appcontext "github.com/turbo514/shortenurl-v2/gateway/app_context"
	"github.com/turbo514/shortenurl-v2/shared/mylog"
	mytrace "github.com/turbo514/shortenurl-v2/shared/trace"
	"github.com/turbo514/shortenurl-v2/shared/util"
	"net/http"
	"strings"
	"time"
)

const (
	USERIDTYPE   = "user_id"
	TENANTIDTYPE = "tenant_id"
)

func AuthMiddleware(app *appcontext.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		tr := mytrace.GetTracer()
		ctx, span := tr.Start(c.Request.Context(), "AuthMiddleware")
		defer span.End()

		authStr := c.GetHeader("Authorization")
		if !strings.HasPrefix(authStr, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"errmsg":  "Authorization header is not Bearer",
				"errcode": 100001,
			})
			return
		}

		token := authStr[7:]
		// TODO: 更新密钥
		claim, err := util.ParseToken(token, nil)
		if err != nil {
			mylog.GetLogger().Warn("令牌解析失败", "err", err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"errmsg":  "无效令牌",
				"errcode": 100002,
			})
			return
		}

		if claim.ExpiresAt.Before(time.Now()) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"errmsg":  "无效令牌",
				"errcode": 100003,
			})
			return
		}

		mylog.GetLogger().Debug("token解析值检查", "user_id", claim.M["user_id"], "tenant_id", claim.M["tenant_id"])

		c.Set("user_id", claim.M["user_id"])
		c.Set("tenant_id", claim.M["tenant_id"])
		c.Request = c.Request.WithContext(ctx)
	}
}
