package middleware

import (
	"bytes"
	"fmt"
	"github.com/MarchGe/go-admin-server/app/common/constant"
	"github.com/gin-gonic/gin"
	"github.com/gobwas/glob"
	"io"
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
		contextPath + "/devops/app/upload",
	}
}

func ApiDebugLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		if ignoredDebug(c.Request.URL.Path) {
			c.Next()
			return
		}
		start := time.Now()
		requestBodyBytes, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewReader(requestBodyBytes))
		c.Next()
		end := time.Now()
		delay := end.Sub(start)
		slog.Debug("RestApiInOutParameters",
			slog.String("requestId", c.GetString(constant.RequestId)),
			slog.String("clientIp", c.ClientIP()),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Any("query", c.Request.URL.Query()),
			slog.String("requestBody", string(requestBodyBytes)),
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
