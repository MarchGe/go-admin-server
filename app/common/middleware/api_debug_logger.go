package middleware

import (
	"fmt"
	"github.com/MarchGe/go-admin-server/app/common/constant"
	"github.com/MarchGe/go-admin-server/app/common/middleware/recorder"
	"github.com/gin-gonic/gin"
	"github.com/gobwas/glob"
	"log/slog"
	"path"
	"time"
)

var ignoreDebugPatterns []string

func initDebugPatterns(contextPath string) {
	ignoreDebugPatterns = []string{
		contextPath + constant.Swagger + "/**",
		contextPath + "/terminal/ws",
		contextPath + "/terminal/ws/ssh/*",
	}
}

func ApiDebugLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		if ignoredDebug(c.Request.URL.Path) {
			slog.Debug("==> Request info: ", slog.String("url", c.Request.URL.String()))
			c.Next()
			return
		}

		start := time.Now()
		requestBodyBytes, bodyIgnored := recorder.GetBodyContent(c)
		c.Next()
		end := time.Now()
		delay := end.Sub(start)
		var bodyLogAttr slog.Attr
		if bodyIgnored {
			bodyLogAttr = slog.Bool("bodyIgnored", true)
		} else {
			bodyLogAttr = slog.String("requestBody", string(requestBodyBytes))
		}

		slog.Debug("==> Request info: ",
			slog.String("requestId", c.GetString(constant.RequestId)),
			slog.String("clientIp", c.ClientIP()),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Any("query", c.Request.URL.Query()),
			bodyLogAttr,
			slog.Duration("duration", delay),
		)
	}
}

func ignoredDebug(requestPath string) bool {
	for _, pattern := range ignoreDebugPatterns {
		g, err := glob.Compile(pattern, '/')
		if err != nil {
			panic(fmt.Errorf("pattern '%s' compile error: %w", pattern, err))
		}
		if g.Match(path.Clean(requestPath)) {
			return true
		}
	}
	return false
}
