package middleware

import (
	"fmt"
	"github.com/MarchGe/go-admin-server/app/common/R"
	"github.com/MarchGe/go-admin-server/app/common/constant"
	"github.com/MarchGe/go-admin-server/app/common/utils"
	"github.com/MarchGe/go-admin-server/config"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gobwas/glob"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
)

var ignorePatterns []string
var websocketPatterns []string

func Initialize(contextPath string) {
	ignorePatterns = []string{
		contextPath + "/auth/login",
		contextPath + "/swagger/**",
	}
	websocketPatterns = []string{
		contextPath + "/terminal/ws",
		contextPath + "/terminal/ws/ssh/*",
	}
	initDebugPatterns(contextPath)
}

func AuthenticationHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if ignored(c.Request.URL.Path) {
			c.Next()
		} else {
			isWs := isWebsocket(c.Request.URL.Path)
			if isWs {
				token := c.Query("token")
				decryptToken, _ := utils.DecryptString(config.GetConfig().Cookie.SecretKey, token, "")
				index := strings.LastIndex(decryptToken, "$")
				if index != -1 {
					loginToken, timestampStr := decryptToken[:index], decryptToken[index+1:]
					timestamp, _ := strconv.Atoi(timestampStr)
					if time.Now().UnixMilli()-int64(timestamp) < 15*1000 {
						c.Request.AddCookie(&http.Cookie{
							Name:  constant.LoginSession,
							Value: loginToken,
						})
					}
				}
			}
			session := sessions.DefaultMany(c, constant.LoginSession)
			userId := session.Get(constant.SessionUserId)
			if userId == nil {
				if isWs {
					c.AbortWithStatus(http.StatusUnauthorized)
				} else {
					R.Fail(c, "授权无效，请重新登录！", http.StatusUnauthorized)
				}
				return
			}
			c.Set(constant.SessionUserId, userId)
			c.Set(constant.IsRootUser, session.Get(constant.IsRootUser))
			c.Next()
		}
	}
}

func ignored(requestPath string) bool {
	for _, pattern := range ignorePatterns {
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

func isWebsocket(requestPath string) bool {
	for _, pattern := range websocketPatterns {
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
