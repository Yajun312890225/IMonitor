package handler

import (
	"fmt"
	"iMonitor/dao"
	"iMonitor/response"
	"iMonitor/utils"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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
	if data.ChildHost == "" {
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
		return
	}
	url := "http://" + data.Host + ":" + data.Port + "/childPing"
	timestamp, sign := utils.GetSign(data.Key1, data.Key2)

	header := req.Header{
		"timestamp": timestamp,
		"sign":      sign,
	}
	param := req.Param{
		"host": data.ChildHost,
		"port": data.ChildPort,
	}
	r, err := req.Post(url, header, param)
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
// @Param parentId query int false "parentId"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.PageResponse "{"code": 200, "data": [...]}"
// @Router /api/v1/serverlist [get]
func GetServerList(c *gin.Context) {
	data := dao.Server()
	var pageSize = 10
	var pageIndex = 1
	var parentId = 0
	if size := c.Request.FormValue("pageSize"); size != "" {
		pageSize, _ = strconv.Atoi(size)
	}

	if index := c.Request.FormValue("pageIndex"); index != "" {
		pageIndex, _ = strconv.Atoi(index)
	}
	if Id := c.Request.FormValue("parentId"); Id != "" {
		parentId, _ = strconv.Atoi(Id)
	}
	session := sessions.Default(c)
	data.CreateBy = strconv.Itoa(session.Get("userid").(int))

	data.Host = c.Request.FormValue("host")
	data.Name = c.Request.FormValue("name")
	data.ParentId = parentId
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
	serverId := func(keys string) (IDS []int) {
		ids := strings.Split(keys, ",")
		for i := 0; i < len(ids); i++ {
			ID, _ := strconv.Atoi(ids[i])
			IDS = append(IDS, ID)
		}
		return

	}(c.Param("serverId"))

	session := sessions.Default(c)

	for _, Id := range serverId {
		if ok, serverId := data.CheckOwner(Id, session.Get("userid").(int)); ok == false {
			err := "没有服务器ID:" + strconv.Itoa(serverId) + "的权限"
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": response.CodeAccessionNotPermission,
				"msg":  err,
			})
			return
		}
	}
	data.UpdateBy = strconv.Itoa(session.Get("userid").(int))

	err := data.BatchDelete(serverId)
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
	session := sessions.Default(c)
	if ok, serverId := server.CheckPermission(server.ServerId, session.Get("userid").(int)); ok == false {
		err := "没有服务器ID:" + strconv.Itoa(serverId) + "的权限"
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": response.CodeAccessionNotPermission,
			"msg":  err,
		})
		return
	}

	if err := server.Get(); err != nil {
		logrus.Info(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}
	if server.ParentId == 0 {
		url := "http://" + server.Host + ":" + server.Port + "/serviceCurrentInfo"
		timestamp, sign := utils.GetSign(server.Key1, server.Key2)
		header := req.Header{
			"timestamp": timestamp,
			"sign":      sign,
		}
		param := req.Param{
			"host": server.IP,
			"port": server.Port,
		}
		r, err := req.Post(url, header, param)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusOK, response.Res{
				Code:  response.CodeParamErr,
				Msg:   response.CodeErrMsg[response.CodeParamErr],
				Error: err.Error(),
			})
			return
		}
		var dic map[string]interface{}
		r.ToJSON(&dic)
		dic["data"] = server
		sc := dao.ServerCollaborator()
		col, err := sc.GetCollaborator(server.ServerId)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusOK, response.Res{
				Code:  response.CodeParamErr,
				Msg:   response.CodeErrMsg[response.CodeParamErr],
				Error: err.Error(),
			})
			return
		}
		dic["collaborator"] = col
		c.JSON(http.StatusOK, dic)
		return
	}
	parentServer := dao.Server()
	parentServer.ServerId = server.ParentId
	if err := parentServer.Get(); err != nil {
		logrus.Info(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}
	url := "http://" + parentServer.Host + ":" + parentServer.Port + "/serviceCurrentInfo"
	timestamp, sign := utils.GetSign(parentServer.Key1, parentServer.Key2)
	header := req.Header{
		"timestamp": timestamp,
		"sign":      sign,
	}
	param := req.Param{
		"host": server.Host,
		"port": server.Port,
	}
	r, err := req.Post(url, header, param)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}
	var dic map[string]interface{}
	r.ToJSON(&dic)
	dic["data"] = server
	c.JSON(http.StatusOK, dic)
}

// UpdateServer 修改服务器信息
// @Summary 修改服务器信息
// @Description 修改服务器信息
// @Tags Server
// @Accept  application/json
// @Product application/json
// @Param data body dao.ReqUpdateServer true "body"
// @Success 200 {string} string	"{"code": 200, "message": "修改成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "修改失败"}"
// @Router /api/v1/server [put]
func UpdateServer(c *gin.Context) {
	server := dao.ReqUpdateServer{}
	err := c.ShouldBindBodyWith(&server, binding.JSON)
	if err != nil {
		logrus.Info(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}
	data := dao.Server()
	data.ServerId = server.ServerId
	data.Name = server.Name
	data.Key1 = server.Key1
	data.Key2 = server.Key2
	data.Sort = server.Sort

	session := sessions.Default(c)
	data.UpdateBy = strconv.Itoa(session.Get("userid").(int))

	result, err := data.Update(data.ServerId)
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

// FetchContact 获取通讯录树
// @Summary 获取通讯录树
// @Description 获取通讯录树
// @Tags Server
// @Accept  application/json
// @Product application/json
// @Param data body dao.ReqFetchContact true "data"
// @Success 200 {string} string	"{"code": 0, "message": "获取成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "获取失败"}"
// @Router /api/v1/fetchContact [post]
func FetchContact(c *gin.Context) {
	data := dao.ReqFetchContact{}
	err := c.ShouldBindBodyWith(&data, binding.JSON)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}
	server := dao.Server()
	server.ServerId = data.ServerId
	if err := server.Get(); err != nil {
		logrus.Info(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}

	if data.OrgId == "" {
		data.OrgId = server.OrgId
	}
	url := "http://" + server.Host + ":" + server.Port + "/contacts/fetchOrgDeptContact"
	timestamp, sign := utils.GetSign(server.Key1, server.Key2)
	header := req.Header{
		"timestamp": timestamp,
		"sign":      sign,
	}
	param := req.Param{
		"pOrgID":      data.OrgId,
		"pDeptID":     data.DeptId,
		"requestType": data.RequestType,
	}
	r, err := req.Post(url, header, param)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}
	var dic map[string]interface{}
	r.ToJSON(&dic)
	c.JSON(http.StatusOK, dic)
}

// SearchUser 查找用户
// @Summary 查找用户
// @Description 查找用户,此接口查找的用户是IM的用户
// @Tags Server
// @Accept  application/json
// @Product application/json
// @Param data body dao.ReqSearchUser true "data"
// @Success 200 {string} string	"{"code": 0, "message": "获取成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "获取失败"}"
// @Router /api/v1/searchUser [post]
func SearchUser(c *gin.Context) {
	data := dao.ReqSearchUser{}
	err := c.ShouldBindBodyWith(&data, binding.JSON)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}
	server := dao.Server()
	server.ServerId = data.ServerId
	if err := server.Get(); err != nil {
		logrus.Info(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}
	url := "http://" + server.Host + ":" + server.Port + "/contacts/searchUser"
	timestamp, sign := utils.GetSign(server.Key1, server.Key2)
	header := req.Header{
		"timestamp": timestamp,
		"sign":      sign,
	}
	param := req.Param{
		"keyword": data.Keyword,
	}
	r, err := req.Post(url, header, param)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}
	var dic map[string]interface{}
	r.ToJSON(&dic)
	c.JSON(http.StatusOK, dic)
}

// FetchUserGroup 查找用户群组
// @Summary 查找用户群组
// @Description 查找用户群组
// @Tags Server
// @Accept  application/json
// @Product application/json
// @Param data body dao.ReqFetchUserGroup true "data"
// @Success 200 {string} string	"{"code": 0, "message": "获取成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "获取失败"}"
// @Router /api/v1/fetchUserGroup [post]
func FetchUserGroup(c *gin.Context) {
	data := dao.ReqFetchUserGroup{}
	err := c.ShouldBindBodyWith(&data, binding.JSON)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}
	server := dao.Server()
	server.ServerId = data.ServerId
	if err := server.Get(); err != nil {
		logrus.Info(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}

	url := "http://" + server.Host + ":" + server.Port + "/contacts/fetchUserGroup"
	timestamp, sign := utils.GetSign(server.Key1, server.Key2)
	header := req.Header{
		"timestamp": timestamp,
		"sign":      sign,
	}
	param := req.Param{
		"userID": data.UserId,
	}
	r, err := req.Post(url, header, param)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}
	var dic map[string]interface{}
	r.ToJSON(&dic)
	c.JSON(http.StatusOK, dic)
}

// FetchMsgRecord 查找聊天记录
// @Summary 查找聊天记录
// @Description 查找聊天记录
// @Tags Server
// @Accept  application/json
// @Product application/json
// @Param data body dao.ReqQueryMsg true "data"
// @Success 200 {string} string	"{"code": 0, "message": "获取成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "获取失败"}"
// @Router /api/v1/fetchMsgRecord [post]
func FetchMsgRecord(c *gin.Context) {
	data := dao.ReqQueryMsg{}
	err := c.ShouldBindBodyWith(&data, binding.JSON)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}
	server := dao.Server()
	server.ServerId = data.ServerId
	if err := server.Get(); err != nil {
		logrus.Info(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}

	if data.PageSize == 0 {
		data.PageSize = 10
	}
	if data.PageIndex == 0 {
		data.PageIndex = 1
	}
	url := "http://" + server.Host + ":" + server.Port + "/contacts/fetchMsgRecord"
	timestamp, sign := utils.GetSign(server.Key1, server.Key2)
	header := req.Header{
		"timestamp": timestamp,
		"sign":      sign,
	}
	param := req.Param{
		"userID":    data.UserId,
		"targetID":  data.TargetId,
		"pageSize":  data.PageSize,
		"pageIndex": data.PageIndex,
		"chatType":  data.ChatType,
		"type":      data.Type,
		"query":     data.Query,
	}
	r, err := req.Post(url, header, param)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}
	var dic map[string]interface{}
	r.ToJSON(&dic)
	c.JSON(http.StatusOK, dic)
}

// CreateCollaborator 添加协作者
// @Summary 添加协作者
// @Description 添加协作者
// @Tags Server
// @Accept  application/json
// @Product applicatio/json
// @Param data body dao.ReqCollaborator true "data"
// @Success 200 {string} string	"{"code": 0, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/createcollaborator [post]
func CreateCollaborator(c *gin.Context) {
	data := dao.ReqCollaborator{}
	err := c.ShouldBindBodyWith(&data, binding.JSON)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}
	server := dao.Server()
	session := sessions.Default(c)
	if ok, serverId := server.CheckOwner(data.ServerId, session.Get("userid").(int)); ok == false {
		err := "没有服务器ID:" + strconv.Itoa(serverId) + "的权限"
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": response.CodeAccessionNotPermission,
			"msg":  err,
		})
		return
	}

	user := dao.User()
	u, err := user.GetUserByName(data.Username)
	if err != nil {
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeUserNotFound,
			Msg:   response.CodeErrMsg[response.CodeUserNotFound],
			Error: err.Error(),
		})
		return
	}
	sc := dao.ServerCollaborator()
	if err := sc.AddCollaborator(data.ServerId, u); err != nil {
		logrus.Error(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeAddCollaboratorErr,
			Msg:   response.CodeErrMsg[response.CodeAddCollaboratorErr],
			Error: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Res{
		Code: response.CodeSuccess,
	})
}

// RemoveCollaborator 移除协作者
// @Summary 移除协作者
// @Description 移除协作者
// @Tags Server
// @Accept  application/json
// @Product applicatio/json
// @Param data body dao.ReqCollaborator true "data"
// @Success 200 {string} string	"{"code": 0, "message": "移除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "移除失败"}"
// @Router /api/v1/removecollaborator [delete]
func RemoveCollaborator(c *gin.Context) {
	data := dao.ReqCollaborator{}
	err := c.ShouldBindBodyWith(&data, binding.JSON)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}
	server := dao.Server()
	session := sessions.Default(c)
	if ok, serverId := server.CheckOwner(data.ServerId, session.Get("userid").(int)); ok == false {
		err := "没有服务器ID:" + strconv.Itoa(serverId) + "的权限"
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": response.CodeAccessionNotPermission,
			"msg":  err,
		})
		return
	}

	user := dao.User()
	u, err := user.GetUserByName(data.Username)
	if err != nil {
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeUserNotFound,
			Msg:   response.CodeErrMsg[response.CodeUserNotFound],
			Error: err.Error(),
		})
		return
	}
	sc := dao.ServerCollaborator()

	if err := sc.DelCollaborator(data.ServerId, u); err != nil {
		logrus.Error(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeDelCollaboratorErr,
			Msg:   response.CodeErrMsg[response.CodeDelCollaboratorErr],
			Error: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Res{
		Code: response.CodeSuccess,
	})
}

// QuerySyncOrgId 获取同步服务器组织
// @Summary 获取同步服务器组织
// @Description 获取同步服务器组织
// @Tags Server
// @Param serverId path string false "serverId"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/querysync/{serverId} [get]
func QuerySyncOrgId(c *gin.Context) {
	server := dao.Server()
	server.ServerId, _ = strconv.Atoi(c.Param("serverId"))
	session := sessions.Default(c)
	if ok, serverId := server.CheckPermission(server.ServerId, session.Get("userid").(int)); ok == false {
		err := "没有服务器ID:" + strconv.Itoa(serverId) + "的权限"
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": response.CodeAccessionNotPermission,
			"msg":  err,
		})
		return
	}

	if err := server.Get(); err != nil {
		logrus.Info(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}

	url := "http://" + server.Host + ":" + server.Port + "/contacts/syncURL"
	timestamp, sign := utils.GetSign(server.Key1, server.Key2)
	header := req.Header{
		"timestamp": timestamp,
		"sign":      sign,
	}
	r, err := req.Post(url, header)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}
	var dic map[string]interface{}
	r.ToJSON(&dic)
	c.JSON(http.StatusOK, dic)
}

// SyncContacts 从OA端同步数据
// @Summary 从OA端同步数据
// @Description 从OA端同步数据
// @Tags Server
// @Accept  application/json
// @Product application/json
// @Param data body dao.ReqSyncContacts true "data"
// @Success 200 {string} string	"{"code": 0, "message": "获取成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "获取失败"}"
// @Router /api/v1/syncContacts [post]
func SyncContacts(c *gin.Context) {
	data := dao.ReqSyncContacts{}
	err := c.ShouldBindBodyWith(&data, binding.JSON)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}

	server := dao.Server()
	server.ServerId = data.ServerId
	if err := server.Get(); err != nil {
		logrus.Info(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}

	url := "http://" + server.Host + ":" + server.Port + "/syncContacts"
	timestamp, sign := utils.GetSign(server.Key1, server.Key2)
	header := req.Header{
		"timestamp": timestamp,
		"sign":      sign,
	}
	param := req.Param{
		"orgId": data.OrgName,
	}
	r, err := req.Post(url, header, param)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}
	var dic map[string]interface{}
	r.ToJSON(&dic)
	if code := dic["code"].(float64); code == 0 {
		res := dic["result"].(map[string]interface{})
		server.OrgId = res["orgId"].(string)
		session := sessions.Default(c)
		server.UpdateBy = strconv.Itoa(session.Get("userid").(int))
		_, err := server.Update(data.ServerId)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusOK, response.Res{
				Code:  response.CodeParamErr,
				Msg:   response.CodeErrMsg[response.CodeParamErr],
				Error: err.Error(),
			})
			return

		}
	}
	c.JSON(http.StatusOK, dic)
}

// UpdateSync 更新同步地址（和添加共用）
// @Summary 更新同步地址（和添加共用）
// @Description 更新同步地址（和添加共用）
// @Tags Server
// @Accept  application/json
// @Product application/json
// @Param data body dao.ReqUpdateSync true "data"
// @Success 200 {string} string	"{"code": 0, "message": "获取成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "获取失败"}"
// @Router /api/v1/updatesync [post]
func UpdateSync(c *gin.Context) {
	data := dao.ReqUpdateSync{}
	err := c.ShouldBindBodyWith(&data, binding.JSON)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}

	server := dao.Server()
	server.ServerId = data.ServerId
	if err := server.Get(); err != nil {
		logrus.Info(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}

	url := "http://" + server.Host + ":" + server.Port + "/contacts/configSyncURL"
	timestamp, sign := utils.GetSign(server.Key1, server.Key2)
	header := req.Header{
		"timestamp": timestamp,
		"sign":      sign,
	}
	param := req.Param{
		"orgID":           data.OrgName,
		"orgSyncUrl":      data.OrgUrl,
		"deptSyncUrl":     data.DeptUrl,
		"userSyncUrl":     data.UserUrl,
		"relationSyncUrl": data.RelationUrl,
	}
	r, err := req.Post(url, header, param)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}
	var dic map[string]interface{}
	r.ToJSON(&dic)
	c.JSON(http.StatusOK, dic)
}

// GetServiceInfo 获取服务器全天信息
// @Summary 获取服务器全天信息
// @Description 获取服务器全天信息
// @Tags Server
// @Accept  application/json
// @Product application/json
// @Param data body dao.ReqServerInfo true "data"
// @Success 200 {string} string	"{"code": 0, "message": "获取成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "获取失败"}"
// @Router /api/v1/serviceInfo [post]
func GetServiceInfo(c *gin.Context) {
	data := dao.ReqServerInfo{}
	err := c.ShouldBindBodyWith(&data, binding.JSON)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}

	server := dao.Server()
	server.ServerId = data.ServerId
	if err := server.Get(); err != nil {
		logrus.Info(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}
	if server.ParentId == 0 {
		url := "http://" + server.Host + ":" + server.Port + "/serviceInfo"
		timestamp, sign := utils.GetSign(server.Key1, server.Key2)
		header := req.Header{
			"timestamp": timestamp,
			"sign":      sign,
		}
		param := req.Param{
			"host": server.IP,
			"port": server.Port,
			"date": data.Date,
		}
		r, err := req.Post(url, header, param)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusOK, response.Res{
				Code:  response.CodeParamErr,
				Msg:   response.CodeErrMsg[response.CodeParamErr],
				Error: err.Error(),
			})
			return
		}
		var dic map[string]interface{}
		r.ToJSON(&dic)
		c.JSON(http.StatusOK, dic)
		return
	}
	parentServer := dao.Server()
	parentServer.ServerId = server.ParentId
	if err := parentServer.Get(); err != nil {
		logrus.Info(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}
	url := "http://" + parentServer.Host + ":" + parentServer.Port + "/serviceInfo"
	timestamp, sign := utils.GetSign(parentServer.Key1, parentServer.Key2)
	header := req.Header{
		"timestamp": timestamp,
		"sign":      sign,
	}
	param := req.Param{
		"host": server.Host,
		"port": server.Port,
		"date": data.Date,
	}
	r, err := req.Post(url, header, param)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}
	var dic map[string]interface{}
	r.ToJSON(&dic)
	c.JSON(http.StatusOK, dic)
}

// UploadServerFile 上传服务器文件
// @Summary 上传服务器文件
// @Description 上传服务器文件
// @Tags Server
// @Accept multipart/form-data
// @Param file formData file true "file"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/uploadfile [post]
func UploadServerFile(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["file"]
	filPath := "static/uploadfile/server/start"
	for _, file := range files {
		// log.Println(file.Filename)
		// 上传文件至指定目录
		timeStr := time.Now().Format("20060102150405")
		_ = os.Rename(filPath, filPath+"_bak_"+timeStr)
		_ = c.SaveUploadedFile(file, filPath)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
	})
}

// RestartServer 重启服务器
// @Summary 重启服务器
// @Description 重启服务器
// @Tags Server
// @Accept  application/json
// @Product application/json
// @Param serverId path string false "serverId"
// @Success 200 {string} string	"{"code": 0, "message": "获取成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "获取失败"}"
// @Router /api/v1/restartserver/{serverId} [get]
func RestartServer(c *gin.Context) {
	server := dao.Server()
	server.ServerId, _ = strconv.Atoi(c.Param("serverId"))
	session := sessions.Default(c)
	if ok, serverId := server.CheckPermission(server.ServerId, session.Get("userid").(int)); ok == false {
		err := "没有服务器ID:" + strconv.Itoa(serverId) + "的权限"
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": response.CodeAccessionNotPermission,
			"msg":  err,
		})
		return
	}

	if err := server.Get(); err != nil {
		logrus.Info(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}

	url := "http://" + server.Host + ":" + server.Port + "/getFile"
	timestamp, sign := utils.GetSign(server.Key1, server.Key2)
	header := req.Header{
		"timestamp": timestamp,
		"sign":      sign,
	}
	r, err := req.Post(url, header)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}
	if status := r.Response().StatusCode; status != http.StatusOK {
		c.JSON(status, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}
	time.Sleep(2 * time.Second)
	url = "http://" + server.Host + ":" + server.Port + "/restart"
	_, _ = req.Post(url, header)

	time.Sleep(5 * time.Second)

	url = "http://" + server.Host + ":" + server.Port + "/ping"

	r, err = req.Post(url, header)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodePingErr,
			Msg:   response.CodeErrMsg[response.CodePingErr],
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
	})
}

// QueryClientLog 获取客户端错误日志
// @Summary 获取客户端错误日志
// @Description 获取客户端错误日志
// @Tags Server
// @Param serverId query int false "serverId"
// @Param query query string false "query"
// @Param pageSize query int false "pageSize"
// @Param pageIndex query int false "pageIndex"
// @Success 200 {string} string	"{"code": 200, "message": "获取成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "获取成功"}"
// @Router /api/v1/queryclientlog [get]
func QueryClientLog(c *gin.Context) {

	server := dao.Server()
	server.ServerId, _ = strconv.Atoi(c.Request.FormValue("serverId"))
	session := sessions.Default(c)
	if ok, serverId := server.CheckPermission(server.ServerId, session.Get("userid").(int)); ok == false {
		err := "没有服务器ID:" + strconv.Itoa(serverId) + "的权限"
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": response.CodeAccessionNotPermission,
			"msg":  err,
		})
		return
	}

	if err := server.Get(); err != nil {
		logrus.Info(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}

	var pageSize = 10
	var pageIndex = 1

	if size := c.Request.FormValue("pageSize"); size != "" {
		pageSize, _ = strconv.Atoi(size)
	}

	if index := c.Request.FormValue("pageIndex"); index != "" {
		pageIndex, _ = strconv.Atoi(index)
	}
	query := c.Request.FormValue("query")

	url := "http://" + server.Host + ":" + server.Port + "/queryLog"
	timestamp, sign := utils.GetSign(server.Key1, server.Key2)
	header := req.Header{
		"timestamp": timestamp,
		"sign":      sign,
	}
	param := req.Param{
		"query":     query,
		"pageSize":  pageSize,
		"pageIndex": pageIndex,
	}
	r, err := req.Post(url, header, param)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}
	var dic map[string]interface{}
	r.ToJSON(&dic)
	c.JSON(http.StatusOK, dic)

}

// GetClientVersion 获取客户端SDK版本号
// @Summary 获取客户端SDK版本号
// @Description 获取客户端SDK版本号
// @Tags Server
// @Param serverId path string false "serverId"
// @Success 200 {string} string	"{"code": 200, "message": "获取成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "获取成功"}"
// @Router /api/v1/getClientVersion/{serverId} [get]
func GetClientVersion(c *gin.Context) {
	server := dao.Server()
	server.ServerId, _ = strconv.Atoi(c.Param("serverId"))
	session := sessions.Default(c)
	if ok, serverId := server.CheckPermission(server.ServerId, session.Get("userid").(int)); ok == false {
		err := "没有服务器ID:" + strconv.Itoa(serverId) + "的权限"
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": response.CodeAccessionNotPermission,
			"msg":  err,
		})
		return
	}

	if err := server.Get(); err != nil {
		logrus.Info(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}

	url := "http://" + server.Host + ":" + server.Port + "/getVersion"
	timestamp, sign := utils.GetSign(server.Key1, server.Key2)
	header := req.Header{
		"timestamp": timestamp,
		"sign":      sign,
	}
	r, err := req.Post(url, header)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}
	var dic map[string]interface{}
	r.ToJSON(&dic)
	c.JSON(http.StatusOK, dic)
}

// UpdateVersion 更新客户端SDK版本号
// @Summary 更新客户端SDK版本号
// @Description 更新客户端SDK版本号
// @Tags Server
// @Accept  application/json
// @Product application/json
// @Param data body dao.ReqUpdateVersion true "data"
// @Success 200 {string} string	"{"code": 0, "message": "获取成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "获取失败"}"
// @Router /api/v1/updateVersion [post]
func UpdateVersion(c *gin.Context) {
	data := dao.ReqUpdateVersion{}
	err := c.ShouldBindBodyWith(&data, binding.JSON)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}

	server := dao.Server()
	server.ServerId = data.ServerId
	if err := server.Get(); err != nil {
		logrus.Info(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}

	url := "http://" + server.Host + ":" + server.Port + "/updateVersion"
	timestamp, sign := utils.GetSign(server.Key1, server.Key2)
	header := req.Header{
		"timestamp": timestamp,
		"sign":      sign,
	}
	param := req.Param{
		"version": data.Version,
	}
	r, err := req.Post(url, header, param)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}
	var dic map[string]interface{}
	r.ToJSON(&dic)
	c.JSON(http.StatusOK, dic)
}
