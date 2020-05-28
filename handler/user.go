package handler

import (
	"net/http"
	"strconv"

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
// @Tags User
// @Success 200 {object} response.Res{data=model.User}
// @Router /api/v1/query [get]
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
// @Tags User
// @Param data body dao.ReqLoginUser true "body"
// @Success 200 {object} response.Res{data=model.User}
// @Router /api/v1/login [post]
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
// @Tags User
// @Success 200 {object} response.Res
// @Router /api/v1/logout [post]
func Logout(c *gin.Context) {
	s := sessions.Default(c)
	s.Clear()
	s.Save()
	c.JSON(http.StatusOK, response.Res{
		Code: response.CodeSuccess,
		Msg:  "成功",
	})
}

// GetUserList 列表数据
// @Summary 列表数据
// @Description 获取JSON
// @Tags User
// @Param username query string false "username"
// @Param phone query string false "phone"
// @Param status query string false "status"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/userlist [get]
func GetUserList(c *gin.Context) {
	data := dao.User()
	var pageSize = 10
	var pageIndex = 1

	if size := c.Request.FormValue("pageSize"); size != "" {
		pageSize, _ = strconv.Atoi(size)
	}

	if index := c.Request.FormValue("pageIndex"); index != "" {
		pageIndex, _ = strconv.Atoi(index)
	}

	data.Username = c.Request.FormValue("username")
	data.Phone = c.Request.FormValue("phone")
	data.Status = c.Request.FormValue("status")

	result, count, err := data.GetPage(pageSize, pageIndex)
	if err != nil {
		logrus.Debug(err)
	}
	c.JSON(http.StatusOK, response.PageResponse{
		Code: response.CodeSuccess,
		Data: response.Page{
			List:      result,
			Count:     count,
			PageIndex: pageIndex,
			PageSize:  pageSize,
		},
		Msg: "",
	})
}
