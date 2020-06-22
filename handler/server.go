package handler

import (
	"fmt"
	"iMonitor/dao"
	"iMonitor/response"
	"iMonitor/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
	"github.com/sirupsen/logrus"
)

// Ping 连接测试
// @Summary 连接测试
// @Description 测试服务器的连通性
// @Tags Server
// @Accept  application/json
// @Product application/json
// @Param data body dao.ReqPing true "data"
// @Success 200 {string} string	"{"code": 0, "message": "连接成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "连接超时"}"
// @Router /api/v1/ping [post]
func Ping(c *gin.Context) {
	data := dao.ReqPing{}
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}
	fmt.Println(data)
	url := "http://" + data.Host + ":" + data.Port + "/ping"
	timestamp, sign := utils.GetSign(data.Key1, data.Key2)

	header := req.Header{
		"timestamp": timestamp,
		"sign":      sign,
	}
	r, err := req.Post(url, header)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodePingErr,
			Msg:   response.CodeErrMsg[response.CodePingErr],
			Error: err.Error(),
		})
		return
	}
	var dic map[string]interface{}
	r.ToJSON(&dic)
	c.JSON(http.StatusOK, dic)

}

// GetServerList 服务器列表
// @Summary 服务器列表
// @Description 服务器列表
// @Tags Server
// @Param host query string false "host"
// @Param name query string false "name"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.PageResponse "{"code": 200, "data": [...]}"
// @Router /api/v1/serverlist [get]
func GetServerList(c *gin.Context) {
	data := dao.Server()
	var pageSize = 10
	var pageIndex = 1

	if size := c.Request.FormValue("pageSize"); size != "" {
		pageSize, _ = strconv.Atoi(size)
	}

	if index := c.Request.FormValue("pageIndex"); index != "" {
		pageIndex, _ = strconv.Atoi(index)
	}

	data.Host = c.Request.FormValue("host")
	data.Name = c.Request.FormValue("name")
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

// AddServer 添加服务器
// @Summary 添加服务器
// @Description 添加服务器
// @Tags Server
// @Accept  application/json
// @Product application/json
// @Param data body dao.ServerDao true "data"
// @Success 200 {string} string	"{"code": 0, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/server [post]
func AddServer(c *gin.Context) {
	data := dao.Server()
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

	c.JSON(http.StatusOK, response.Res{
		Code: response.CodeSuccess,
		Data: result,
		Msg:  "",
	})
}

// DeleteServer 删除服务器
// @Summary 删除服务器
// @Description 删除服务器
// @Tags Server
// @Accept  application/json
// @Product application/json
// @Param serverId path string true "serverId"
// @Success 200 {string} string	"{"code": 0, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/server/{serverId} [delete]
func DeleteServer(c *gin.Context) {
	data := dao.Server()
	roleIds := func(keys string) (IDS []int) {
		ids := strings.Split(keys, ",")
		for i := 0; i < len(ids); i++ {
			ID, _ := strconv.Atoi(ids[i])
			IDS = append(IDS, ID)
		}
		return

	}(c.Param("serverId"))

	session := sessions.Default(c)
	data.UpdateBy = strconv.Itoa(session.Get("userid").(int))

	err := data.BatchDelete(roleIds)
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

// GetServer 获取服务器信息
// @Summary 获取服务器信息
// @Description 获取服务器信息
// @Tags Server
// @Param serverId path string false "serverId"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/server/{serverId} [get]
func GetServer(c *gin.Context) {
	server := dao.Server()
	server.ServerId, _ = strconv.Atoi(c.Param("serverId"))
	if err := server.Get(); err != nil {
		logrus.Info(err)
	}
	url := "http://" + server.Host + ":" + server.Port + "/serviceCurrentInfo"
	timestamp, sign := utils.GetSign(server.Key1, server.Key2)
	header := req.Header{
		"timestamp": timestamp,
		"sign":      sign,
	}
	r, err := req.Post(url, header)
	if err != nil {
		logrus.Error(err)
	}

	var dic map[string]interface{}
	r.ToJSON(&dic)
	c.JSON(http.StatusOK, dic)

}
