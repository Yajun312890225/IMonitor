package handler

import (
	"iMonitor/dao"
	"iMonitor/response"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GetMenuList 获取Menu列表
// @Summary Menu列表数据
// @Description 获取JSON
// @Tags Menu
// @Param menuName query string false "menuName"
// @Param visible query string false "visible"
// @Param title query string false "title"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/menulist [get]
func GetMenuList(c *gin.Context) {

	menu := dao.Menu()
	menu.MenuName = c.Request.FormValue("menuName")
	menu.Visible = c.Request.FormValue("visible")
	menu.Title = c.Request.FormValue("title")
	result, err := menu.GetAllMenu()
	if err != nil {
		logrus.Info(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Res{
		Code: response.CodeSuccess,
		Data: result,
		Msg:  "",
	})
}

// InsertMenu 创建菜单
// @Summary 创建菜单
// @Description 获取JSON
// @Tags Menu
// @Accept  application/json
// @Product application/json
// @Param data body model.Menu true "data"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/menu [post]
func InsertMenu(c *gin.Context) {
	data := dao.Menu()
	err := c.ShouldBind(&data)
	if err != nil {
		logrus.Info(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}
	session := sessions.Default(c)
	data.CreateBy = strconv.Itoa(session.Get("userid").(int))
	result, err := data.Create()
	if err != nil {
		logrus.Info(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeMenuCreateErr,
			Msg:   response.CodeErrMsg[response.CodeMenuCreateErr],
			Error: err.Error(),
		})
		return

	}
	c.JSON(http.StatusOK, response.Res{
		Code: response.CodeSuccess,
		Data: result,
		Msg:  "",
	})

}

// UpdateMenu  修改菜单
// @Summary 修改菜单
// @Description 获取JSON
// @Tags Menu
// @Accept  application/json
// @Product application/json
// @Param data body model.Menu true "body"
// @Success 200 {string} string	"{"code": 200, "message": "修改成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "修改失败"}"
// @Router /api/v1/menu [put]
func UpdateMenu(c *gin.Context) {
	data := dao.Menu()
	err := c.ShouldBind(&data)
	logrus.Info(data.MenuId)
	if err != nil {
		logrus.Info(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}
	session := sessions.Default(c)
	data.UpdateBy = strconv.Itoa(session.Get("userid").(int))
	_, err = data.Update(data.MenuId)
	if err != nil {
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeMenuUpdateErr,
			Msg:   response.CodeErrMsg[response.CodeMenuUpdateErr],
			Error: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Res{
		Code: response.CodeSuccess,
		Data: "",
		Msg:  "修改成功",
	})

}

// @Summary 删除菜单
// @Description 删除数据
// @Tags Menu
// @Param id path int true "id"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/menu/{id} [delete]
func DeleteMenu(c *gin.Context) {
	data := dao.Menu()
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		logrus.Info(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}
	session := sessions.Default(c)
	data.UpdateBy = strconv.Itoa(session.Get("userid").(int))
	_, err = data.Delete(id)
	if err != nil {
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Res{
		Code: response.CodeSuccess,
		Data: "",
		Msg:  "删除成功",
	})
}
