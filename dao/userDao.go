package dao

import (
	"errors"
	"iMonitor/model"
	"iMonitor/response"
)

// UserDao 对user模型进行增删查改
type UserDao struct {
	model.User
}

// User
func User() *UserDao {
	return &UserDao{}
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

// 获取用户数据
func (u *UserDao) Get() (user UserDao, err error) {

	table := model.DB.Table("user").Select([]string{"user.*", "role.role_name"})
	table = table.Joins("left join role on user.role_id=role.role_id")
	if u.UserId != 0 {
		table = table.Where("user_id = ?", u.UserId)
	}

	if u.Username != "" {
		table = table.Where("username = ?", u.Username)
	}

	if u.Password != "" {
		table = table.Where("password = ?", u.Password)
	}

	if u.RoleId != 0 {
		table = table.Where("role_id = ?", u.RoleId)
	}

	if err = table.First(&user).Error; err != nil {
		return
	}
	return
}

//Insert 添加用户
func (u UserDao) Insert() (id int, err error) {
	// check 用户名
	var count int
	model.DB.Table("user").Where("username = ?", u.Username).Count(&count)
	if count > 0 {
		err = errors.New("账户已存在！")
		return
	}

	//添加数据
	if err = model.DB.Table("user").Create(&u).Error; err != nil {
		return
	}
	id = u.UserId
	return
}
