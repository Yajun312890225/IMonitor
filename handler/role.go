package handler

import (
	"iMonitor/dao"
	"iMonitor/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GetRoleList 角色列表数据
// @Summary 角色列表数据
// @Description Get JSON
// @Tags Role
// @Param roleName query string false "roleName"
// @Param status query string false "status"
// @Param rolekey query string false "rolekey"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.PageResponse "{"code": 200, "data": [...]}"
// @Router /api/v1/rolelist [get]
// @Security
func GetRoleList(c *gin.Context) {
	var pageSize = 10
	var pageIndex = 1

	if size := c.Request.FormValue("pageSize"); size != "" {
		pageSize, _ = strconv.Atoi(size)
	}

	if index := c.Request.FormValue("pageIndex"); index != "" {
		pageIndex, _ = strconv.Atoi(index)
	}
	result, count, err := dao.Role().GetPage(pageSize, pageIndex)
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

// GetRole 获取Role数据
// @Summary 获取Role数据
// @Description 获取JSON
// @Tags Role
// @Param roleId path string false "roleId"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/role/{roleId} [get]
// @Security
func GetRole(c *gin.Context) {
	role := dao.Role()
	role.RoleId, _ = strconv.Atoi(c.Param("roleId"))
	if err := role.Get(); err != nil {
		logrus.Info(err)
	}
	menuIds := make([]int, 0)
	menuIds, err := role.GetRoleMeunId()
	if err != nil {
		logrus.Info(err)
	}
	role.MenuIds = menuIds
	c.JSON(http.StatusOK, response.Res{
		Code: response.CodeSuccess,
		Data: role,
		Msg:  "",
	})

}
