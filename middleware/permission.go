package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"

	"iMonitor/dao"
	"iMonitor/model"
	"iMonitor/pkg/casbin"
	"iMonitor/response"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

//AuthCheckRole 权限检查中间件
func AuthCheckRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		cas := casbin.GetCasbin()
		//获取请求的URI
		obj := c.Request.URL.Path
		//获取请求方法
		act := c.Request.Method
		//获取用户的角色
		// sub := c.Query("rolekey")
		session := sessions.Default(c)
		sub := session.Get("rolekey")
		//判断策略中是否存在

		if b := cas.Enforce.Enforce(sub, obj, act); !b {
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

func CheckPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		type PreMission struct {
			ServerId int `json:"serverId" binding:"required"`
		}
		data := PreMission{}
		if err := c.ShouldBindBodyWith(&data, binding.JSON); err != nil {
			c.JSON(http.StatusOK, response.Res{
				Code:  response.CodeParamErr,
				Msg:   response.CodeErrMsg[response.CodeParamErr],
				Error: err.Error(),
			})
			c.Abort()
			return
		}

		session := sessions.Default(c)
		dao := dao.Server()
		if ok, serverId := dao.CheckPermission(data.ServerId, session.Get("userid").(int)); ok == false {
			err := "没有服务器ID:" + strconv.Itoa(serverId) + "的权限"
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": response.CodeAccessionNotPermission,
				"msg":  err,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
