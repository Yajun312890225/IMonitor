package model

import "time"

// Model 基础类型
type Model struct {
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

// User 用户
type User struct {
	UserId    int    `gorm:"primary_key;AUTO_INCREMENT"  json:"userId"`
	Nickname  string `gorm:"type:varchar(128)" json:"nickname"`
	Phone     string `gorm:"type:varchar(11)" json:"phone"`
	RoleId    int    `gorm:"type:int(11)" json:"roleId"`
	Username  string `gorm:"type:varchar(64)" json:"username"`
	Password  string `gorm:"type:varchar(128)" json:"password"`
	Salt      string `gorm:"type:varchar(255)" json:"salt"`
	Avatar    string `gorm:"type:varchar(255)" json:"avatar"`
	Sex       string `gorm:"type:varchar(255)" json:"sex"`
	Email     string `gorm:"type:varchar(128)" json:"email"`
	Status    string `gorm:"type:enum('admin', 'normal');default:'normal'" json:"status"`
	CreateBy  string `gorm:"type:varchar(128)" json:"createBy"`
	UpdateBy  string `gorm:"type:varchar(128)" json:"updateBy"`
	Remark    string `gorm:"type:varchar(255)" json:"remark"`
	DataScope string `gorm:"-" json:"dataScope"`
	Params    string `gorm:"-" json:"params"`
	Model
}

//Role 角色
type Role struct {
	RoleId    int    `json:"roleId" gorm:"primary_key;AUTO_INCREMENT"`
	RoleName  string `json:"roleName" gorm:"type:varchar(128);"`
	Status    string `json:"status" gorm:"type:int(1);"`
	RoleKey   string `json:"roleKey" gorm:"type:varchar(128);"`
	RoleSort  int    `json:"roleSort" gorm:"type:int(4);"`
	DataScope string `json:"dataScope" gorm:"type:varchar(128);"`
	CreateBy  string `json:"createBy" gorm:"type:varchar(128);"`
	UpdateBy  string `json:"updateBy" gorm:"type:varchar(128);"`
	Remark    string `json:"remark" gorm:"type:varchar(255);"`
	Params    string `json:"params" gorm:"-"`
	MenuIds   []int  `json:"menuIds" gorm:"-"`
	Model
}

type Menu struct {
	MenuId     int    `form:"menuId" json:"menuId" gorm:"primary_key;AUTO_INCREMENT"`
	MenuName   string `form:"menuName" json:"menuName" gorm:"type:varchar(11);"`
	Title      string `form:"title" json:"title" gorm:"type:varchar(64);"`
	Icon       string `form:"icon" json:"icon" gorm:"type:varchar(128);"`
	Path       string `form:"path" json:"path" gorm:"type:varchar(128);"`
	Paths      string `form:"paths" json:"paths" gorm:"type:varchar(128);"`
	MenuType   string `form:"menuType" json:"menuType" gorm:"type:varchar(1);"`
	Action     string `form:"action" json:"action" gorm:"type:varchar(16);"`
	Permission string `form:"permission" json:"permission" gorm:"type:varchar(32);"`
	ParentId   int    `form:"parentId" json:"parentId" gorm:"type:int(11);"`
	NoCache    bool   `form:"noCache" json:"noCache" gorm:"type:char(1);"`
	Breadcrumb string `form:"breadcrumb" json:"breadcrumb" gorm:"type:varchar(255);"`
	Component  string `form:"component" json:"component" gorm:"type:varchar(255);"`
	Sort       int    `form:"sort" json:"sort" gorm:"type:int(4);"`
	Visible    string `form:"visible" json:"visible" gorm:"type:char(1);"`
	CreateBy   string `json:"createBy" gorm:"type:varchar(128);"`
	UpdateBy   string `json:"updateBy" gorm:"type:varchar(128);"`
	IsFrame    string `json:"isFrame" gorm:"type:int(1);DEFAULT:0;"`
	Params     string `json:"params" gorm:"-"`
	RoleId     int    `gorm:"-"`
	Children   []Menu `json:"children" gorm:"-"`
	IsSelect   bool   `json:"is_select" gorm:"-"`
	Model
}

//RoleMenu 角色菜单绑定
type RoleMenu struct {
	RoleId   int    `gorm:"type:int(11)"`
	MenuId   int    `gorm:"type:int(11)"`
	RoleName string `gorm:"type:varchar(128)"`
}

// Server IM服务器
type Server struct {
	Model
	Host    string `json:"host"`
	Port    string `json:"port"`
	Name    string `json:"name"`
	Key1    string `json:"key1"`
	Key2    string `json:"key2"`
	Manager uint   `json:"manager"`
}

// ServerCollaborator 服务器协作者
type ServerCollaborator struct {
	Serverid uint
	Userid   uint
}
