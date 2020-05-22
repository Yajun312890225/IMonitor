package handler

import (
	"iMonitor/dao"
	"iMonitor/response"
	"net/http"

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
	}
	c.JSON(http.StatusOK, response.Res{
		Code: response.CodeSuccess,
		Data: result,
		Msg:  "",
	})
}
