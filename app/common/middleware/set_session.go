package middleware

import (
	"github.com/MarchGe/go-admin-server/app/common/constant"
	"github.com/MarchGe/go-admin-server/config"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func SetSession(cookieConfig *config.CookieConfig) gin.HandlerFunc {
	store := createSessionStore(cookieConfig)
	return sessions.SessionsMany([]string{constant.LoginSession}, store)
}

func createSessionStore(cookieCfg *config.CookieConfig) sessions.Store {
	authenticationKey := []byte(cookieCfg.AuthenticationKey)
	secretKey := []byte(cookieCfg.SecretKey)
	store := cookie.NewStore(authenticationKey, secretKey)
	store.Options(sessions.Options{
		Path:     cookieCfg.Path,
		MaxAge:   cookieCfg.MaxAge,
		Secure:   cookieCfg.Secure,
		HttpOnly: cookieCfg.HttpOnly,
	})
	return store
}
