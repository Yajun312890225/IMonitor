package dao

import (
	"iMonitor/model"
	"iMonitor/response"
)

// UserDao 对user模型进行增删查改的单例工具
type UserDao struct {
	model.User
}

var userDao *UserDao

// User 得到dao-user 单例工具
func User() *UserDao {
	if userDao == nil {
		userDao = &UserDao{}
	}
	return userDao
}

// ReqLoginUser 用户登录用来解析账号密码
type ReqLoginUser struct {
	Username string `form:"username" json:"username" xml:"username" binding:"required" example:"admin"`
	Password string `form:"password" json:"password" xml:"password" binding:"required" example:"123"`
}

// ReqAddUser 管理员新增用户
type ReqAddUser struct {
	Username string `form:"username" binding:"required"`
	Nickname string `form:"nickname" binding:"required"`
}

// Login 去数据库验证登录
func (reqLoginUser *ReqLoginUser) Login(block func(*model.User)) response.Res {
	var user model.User
	if err := model.DB.Where("username = ?", reqLoginUser.Username).First(&user).Error; err != nil {
		return response.Res{
			Code: response.CodeUserNotFound,
			Msg:  response.CodeErrMsg[response.CodeUserNotFound],
		}
	}
	if user.Password != reqLoginUser.Password {
		return response.Res{
			Code: response.CodePasswordErr,
			Msg:  response.CodeErrMsg[response.CodePasswordErr],
		}
	}
	// 登录成功，清楚之前储存的userId，重新设置userId
	block(&user)

	return response.Res{
		Code: response.CodeSuccess,
		Data: user,
	}
}

// RegistUser 添加新用户
func (addUser *ReqAddUser) RegistUser() response.Res {
	var user = model.User{
		Username: addUser.Username,
		Password: "111111", // 默认密码
		Nickname: addUser.Nickname,
		Status:   "1", // 默认普通用户
	}
	if err := model.DB.Create(&user).Error; err != nil {
		return response.Res{
			Code:  response.CodeParamErr,
			Msg:   response.CodeErrMsg[response.CodeParamErr],
			Error: err.Error(),
		}
	}
	return response.Res{
		Code: response.CodeSuccess,
		Data: user,
	}
}

// GetUserByID 通过id查询user
func (*UserDao) GetUserByID(id interface{}) (model.User, error) {
	var uesr model.User
	result := model.DB.First(&uesr, id)
	return uesr, result.Error
}

// GetPage 获取用户列表
func (u *UserDao) GetPage(pageSize int, pageIndex int) ([]UserDao, int, error) {
	var doc []UserDao
	table := model.DB.Select("user.*").Table("user")
	if u.Username != "" {
		table = table.Where("username = ?", u.Username)
	}
	if u.Status != "" {
		table = table.Where("user.status = ?", u.Status)
	}

	if u.Phone != "" {
		table = table.Where("user.phone = ?", u.Phone)
	}

	var count int
	if err := table.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&doc).Error; err != nil {
		return nil, 0, err
	}
	table.Where("user.deleted_at IS NULL").Count(&count)
	return doc, count, nil
}
