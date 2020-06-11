package handler

import (
	"iMonitor/dao"
	"iMonitor/response"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// GetInfo 获取权限信息
// @Summary  获取权限信息
// @Description 获取JSON
// @Tags Info
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/getinfo [get]
func GetInfo(c *gin.Context) {

	var roles = make([]string, 1)
	session := sessions.Default(c)
	roles[0] = session.Get("rolekey").(string)

	var permissions = make([]string, 1)
	permissions[0] = "*:*:*"
	RoleMenu := dao.RoleMenu()
	RoleMenu.RoleId = session.Get("roleId").(int)

	var mp = make(map[string]interface{})
	mp["roles"] = roles
	// if session.Get("rolekey").(string) == "admin" {
	// 	mp["permissions"] = permissions
	// } else {

	// }
	list, _ := RoleMenu.GetPermis()
	mp["permissions"] = list

	c.JSON(http.StatusOK, response.Res{
		Code: response.CodeSuccess,
		Data: mp,
		Msg:  "",
	})
}
