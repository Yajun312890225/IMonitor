package handler

import (
	"iMonitor/dao"
	"iMonitor/model"
	"iMonitor/response"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-contrib/sessions"
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
func GetRoleList(c *gin.Context) {
	data := dao.Role()
	var pageSize = 10
	var pageIndex = 1

	if size := c.Request.FormValue("pageSize"); size != "" {
		pageSize, _ = strconv.Atoi(size)
	}

	if index := c.Request.FormValue("pageIndex"); index != "" {
		pageIndex, _ = strconv.Atoi(index)
	}
	data.RoleKey = c.Request.FormValue("roleKey")
	data.RoleName = c.Request.FormValue("roleName")
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

// GetRole 获取Role数据
// @Summary 获取Role数据
// @Description 获取JSON
// @Tags Role
// @Param roleId path string false "roleId"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/role/{roleId} [get]
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

// InsertRole 创建角色
// @Summary 创建角色
// @Description 获取JSON
// @Tags Role
// @Accept  application/json
// @Product application/json
// @Param data body model.Role true "data"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/role [post]
func InsertRole(c *gin.Context) {

	data := dao.Role()
	role := model.Role{}
	err := c.ShouldBind(&role)
	data.Role = role
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

	result, err := data.Insert()
	if err != nil {
		logrus.Info(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeRoleCreateErr,
			Msg:   response.CodeErrMsg[response.CodeRoleCreateErr],
			Error: err.Error(),
		})
		return

	}
	data.RoleId = result
	roleMenu := dao.RoleMenu()
	err = roleMenu.Insert(result, data.MenuIds)
	if err != nil {
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeRoleMenuCreateErr,
			Msg:   response.CodeErrMsg[response.CodeRoleMenuCreateErr],
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.Res{
		Code: response.CodeSuccess,
		Data: data,
		Msg:  "",
	})
}

// UpdateRole 修改用户角色
// @Summary 修改用户角色
// @Description 获取JSON
// @Tags Role
// @Accept  application/json
// @Product application/json
// @Param data body model.Role true "body"
// @Success 200 {string} string	"{"code": 200, "message": "修改成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "修改失败"}"
// @Router /api/v1/role [put]
func UpdateRole(c *gin.Context) {
	data := dao.Role()
	role := model.Role{}
	err := c.ShouldBind(&role)
	data.Role = role
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

	result, err := data.Update(data.RoleId)
	if err != nil {
		logrus.Info(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeRoleUpdateErr,
			Msg:   response.CodeErrMsg[response.CodeRoleUpdateErr],
			Error: err.Error(),
		})
		return

	}
	roleMenu := dao.RoleMenu()
	err = roleMenu.DeleteRoleMenu(data.RoleId)
	if err != nil {
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeRoleMenuUpdateErr,
			Msg:   response.CodeErrMsg[response.CodeRoleMenuUpdateErr],
			Error: err.Error(),
		})
		return
	}
	err = roleMenu.Insert(data.RoleId, data.MenuIds)
	if err != nil {
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeRoleMenuUpdateErr,
			Msg:   response.CodeErrMsg[response.CodeRoleMenuUpdateErr],
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

// DeleteRole 删除用户角色
// @Summary 删除用户角色
// @Description 删除数据
// @Tags Role
// @Param roleId path string true "roleId"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/role/{roleId} [delete]
func DeleteRole(c *gin.Context) {

	data := dao.Role()
	roleIds := func(keys string) (IDS []int) {
		ids := strings.Split(keys, ",")
		for i := 0; i < len(ids); i++ {
			ID, _ := strconv.Atoi(ids[i])
			IDS = append(IDS, ID)
		}
		return

	}(c.Param("roleId"))

	session := sessions.Default(c)
	data.UpdateBy = strconv.Itoa(session.Get("userid").(int))

	roleMenu := dao.RoleMenu()
	err := roleMenu.BatchDeleteRoleMenu(roleIds)
	if err != nil {
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeRoleMenuUpdateErr,
			Msg:   response.CodeErrMsg[response.CodeRoleMenuUpdateErr],
			Error: err.Error(),
		})
		return
	}

	err = data.BatchDelete(roleIds)
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
