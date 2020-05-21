package handler

import (
	"net/http"

	"iMonitor/dao"
	"iMonitor/model"
	"iMonitor/response"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Query 客户端可以用此链接查询自身登录状态。cookie没过期，返回user
// @Summary 查询登录状态，并返回user
// @Description 查询登录状态，并返回user
// @Tags user
// @Success 200 {object} response.Res{data=model.User}
// @Router /query [post]
func Query(c *gin.Context) {
	if user, _ := c.Get("user"); user != nil {
		if _, ok := user.(*model.User); ok {
			c.JSON(http.StatusOK, response.Res{
				Code: response.CodeSuccess,
				Data: user,
			})
		} else {
			c.JSON(http.StatusOK, response.Res{
				Code: response.CodeCheckLogin,
				Msg:  response.CodeErrMsg[response.CodeCheckLogin],
			})
		}
	} else {
		c.JSON(http.StatusOK, response.Res{
			Code: response.CodeCheckLogin,
			Msg:  response.CodeErrMsg[response.CodeCheckLogin],
		})
	}
}

// Login 登录
// @Summary 用户登录
// @Tags user
// @Param username formData string true "admin"
// @Param password formData string true "admin@123"
// @Success 200 {object} response.Res{data=model.User}
// @Router /login [post]
func Login(c *gin.Context) {
	var loginDao dao.ReqLoginUser
	if err := c.ShouldBind(&loginDao); err == nil {
		res := loginDao.Login(func(user *model.User) {
			session := sessions.Default(c)
			session.Clear()
			session.Set("userid", user.UserId)
			role := dao.Role()
			role.RoleId = user.RoleId
			if err := role.Get(); err != nil {
				logrus.Info(err)
			}
			session.Set("rolekey", role.RoleKey)

			session.Save()
		})
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
	}
}

// Logout 登出
// @Summary 注销登录
// @Tags user
// @Success 200 {object} response.Res
// @Router /user/logout [post]
func Logout(c *gin.Context) {
	s := sessions.Default(c)
	s.Clear()
	s.Save()
	c.JSON(http.StatusOK, response.Res{
		Code: response.CodeSuccess,
		Msg:  "成功",
	})
}

// AddUser 管理员添加用户
// @Summary 管理员添加用户
// @Tags user
// @Param username formData string true "dahuang"
// @Param nickname formData string true "大黄"
// @Success 200 {object} response.Res{data=model.User}
// @Router /user/add [post]
func AddUser(c *gin.Context) {
	var addDao dao.ReqAddUser
	if err := c.ShouldBind(&addDao); err != nil {
		res := addDao.RegistUser()
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
	}
}

func DeleteUser(c *gin.Context) {
	// var userid string
}

func UserList(c *gin.Context) {

}
