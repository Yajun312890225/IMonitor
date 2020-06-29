package handler

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"iMonitor/dao"
	"iMonitor/model"
	"iMonitor/response"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
			session.Set("roleId", role.RoleId)
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

// GetUser 获取用户
// @Summary 获取用户
// @Description 获取JSON
// @Tags User
// @Param userId path int true "userId"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/user/{userId} [get]
func GetUser(c *gin.Context) {
	data := dao.User()
	data.UserId, _ = strconv.Atoi(c.Param("userId"))
	result, err := data.Get()
	if err != nil {
		logrus.Debug(err)
	}
	roles, err := dao.Role().GetList()

	roleIds := make([]int, 0)
	roleIds = append(roleIds, result.RoleId)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    result,
		"roleIds": roleIds,
		"roles":   roles,
	})
}

// InsertUser 创建用户
// @Summary 创建用户
// @Description 获取JSON
// @Tags User
// @Accept  application/json
// @Product application/json
// @Param data body dao.ReqAddUser true "用户数据"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/user [post]
func InsertUser(c *gin.Context) {

	user := dao.ReqAddUser{}
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}
	session := sessions.Default(c)
	user.CreateBy = strconv.Itoa(session.Get("userid").(int))
	id, err := user.Insert()
	if err != nil {
		logrus.Debug(err)
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Res{
		Code: response.CodeSuccess,
		Data: id,
		Msg:  "",
	})

}

// UpdateUser 修改用户数据
// @Summary 修改用户数据
// @Description 获取JSON
// @Tags User
// @Accept  application/json
// @Product application/json
// @Param data body dao.ReqUpdateUser true "body"
// @Success 200 {string} string	"{"code": 200, "message": "修改成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "修改失败"}"
// @Router /api/v1/user [put]
func UpdateUser(c *gin.Context) {
	// data := dao.User()
	// user := model.User{}
	user := dao.ReqUpdateUser{}
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusOK, response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		})
		return
	}
	// data.User = user
	session := sessions.Default(c)
	user.UpdateBy = strconv.Itoa(session.Get("userid").(int))
	result, err := user.Update(user.UserId)
	if err != nil {
		logrus.Debug(err)
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

// DeleteUser 删除用户数据
// @Summary 删除用户数据
// @Description 删除数据
// @Tags User
// @Param userId path int true "userId"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/user/{userId} [delete]
func DeleteUser(c *gin.Context) {
	data := dao.User()
	userId := func(keys string) (IDS []int) {
		ids := strings.Split(keys, ",")
		for i := 0; i < len(ids); i++ {
			ID, _ := strconv.Atoi(ids[i])
			IDS = append(IDS, ID)
		}
		return
	}(c.Param("userId"))

	session := sessions.Default(c)
	data.UpdateBy = strconv.Itoa(session.Get("userid").(int))
	result, err := data.BatchDelete(userId)
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
		Data: result,
		Msg:  "删除成功",
	})
}

// @Summary 修改头像
// @Description 获取JSON
// @Tags User
// @Accept multipart/form-data
// @Param file formData file true "file"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/user/profileAvatar [post]
func InsertUserAvatar(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["file"]
	guid := uuid.New().String()
	filPath := "static/uploadfile/" + guid + ".jpg"
	for _, file := range files {
		log.Println(file.Filename)
		// 上传文件至指定目录
		_ = c.SaveUploadedFile(file, filPath)
	}
	user := dao.ReqUpdateUser{}
	session := sessions.Default(c)
	user.UserId = session.Get("userid").(int)
	user.Avatar = os.Getenv("UPLOADFILE") + "/" + filPath
	user.UpdateBy = strconv.Itoa(session.Get("userid").(int))
	result, err := user.Update(user.UserId)
	if err != nil {
		logrus.Debug(err)
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

// GetAllUser 获取所有用户
// @Summary 获取所有用户
// @Description 获取所有用户
// @Tags User
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/alluser [get]
func GetAllUser(c *gin.Context) {
	data := dao.User()
	result, err := data.GetAllUser()
	if err != nil {
		logrus.Debug(err)
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
