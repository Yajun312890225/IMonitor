package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/sirupsen/logrus"

	"iMonitor/model"
	"iMonitor/pkg/casbin"
	"iMonitor/response"

	"github.com/gin-gonic/gin"
)

//AuthCheckRole 权限检查中间件
func AuthCheckRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		e, err := casbin.Casbin()
		if err != nil {
			logrus.Debug(err)
		}

		//获取请求的URI
		obj := c.Request.URL.Path
		//获取请求方法
		act := c.Request.Method
		//获取用户的角色
		// sub := c.Query("rolekey")
		session := sessions.Default(c)
		sub := session.Get("rolekey")
		//判断策略中是否存在

		if b := e.Enforce(sub, obj, act); !b {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": response.CodeAccessionNotPermission,
				"msg":  response.CodeErrMsg[response.CodeAccessionNotPermission],
			})
			c.Abort()
			return
		} else {
			if user, _ := c.Get("user"); user != nil {
				if _, ok := user.(*model.User); ok {
					c.Next()
					return
				}
			}
		}
	}
}
