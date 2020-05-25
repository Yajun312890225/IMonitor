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
// @Security Bearer
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
// @Accept  application/x-www-form-urlencoded
// @Product application/x-www-form-urlencoded
// @Param menuName formData string true "menuName"
// @Param title formData string true "title"
// @Param menuType formData string true "menuType"
// @Param path formData string true "path"
// @Param permission formData string true "permission"
// @Param action formData string false "action"
// @Param parentId formData string false "parentId"
// @Param isFrame formData string false "isFrame"
// @Param sort formData string false "sort"
// @Param visible formData string false "visible"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/menu [post]
// @Security Bearer
func InsertMenu(c *gin.Context) {
	data := dao.Menu()
	// var addmenu dao.ReqInsertMenu
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
