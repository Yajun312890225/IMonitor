package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// Session 添加中间件操作 session --- s := sessions.Default(c *gin.Context)
func Session(secret string) gin.HandlerFunc {
	store := cookie.NewStore([]byte(secret))
	store.Options(sessions.Options{
		HttpOnly: true,
		MaxAge:   7 * 86400,
		Path:     "/",
	})
	return sessions.Sessions("imonitor-session", store)
}
