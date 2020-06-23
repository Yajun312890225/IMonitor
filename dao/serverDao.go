package dao

import (
	"iMonitor/model"
)

type ServerDao model.Server
type ReqPing struct {
	Host string `json:"host" binding:"required"`
	Port string `json:"port" binding:"required"`
	Key1 string `json:"key1" binding:"required"`
	Key2 string `json:"key2" binding:"required"`
}
type ReqUpdateServer struct {
	ServerId int    `json:"serverId" binding:"required"`
	Name     string `json:"name"`
	Key1     string `json:"key1"`
	Key2     string `json:"key2"`
	OrgId    string `json:"orgId"`
	Sort     int    `json:"sort"`
}

type ReqFetchContact struct {
	ServerId    int    `json:"serverId" binding:"required"`
	OrgId       string `json:"orgId"`       // 父组织id，最顶层需要用户自己输入
	DeptId      string `json:"deptId"`      // 父部门id，组织下的直系部门为ROOT
	RequestType int    `json:"requestType"` // 1：组织 2：部门
}
type ReqSearchUser struct {
	ServerId int    `json:"serverId" binding:"required"`
	Keyword  string `json:"keyword" binding:"required"` //模糊查询
}
type ReqFetchUserGroup struct {
	ServerId int    `json:"serverId" binding:"required"`
	UserId   string `json:"userId" binding:"required"`
}

type ReqQueryMsg struct {
	ServerId  int    `json:"serverId" binding:"required"`
	UserId    string `json:"userId" binding:"required"`   // 本人userID
	TargetId  string `json:"targetId" binding:"required"` // 目标userID
	PageSize  int    `json:"pageSize" `                   // 每页容量，默认10
	PageIndex int    `json:"pageIndex" `                  // 页索引 从1开始
	ChatType  int    `json:"chatType"`                    // 0:单聊 ,1:群聊
	Type      int    `json:"type"`                        // 消息类型，目前预留
	Query     string `json:"query"`                       // 搜索关键字
}

func Server() *ServerDao {
	return &ServerDao{}
}

//GetPage 获取服务器
func (s *ServerDao) GetPage(pageSize int, pageIndex int) ([]model.Server, int, error) {
	var doc []model.Server

	table := model.DB.Select("server.*").Table("server").Joins("LEFT JOIN server_collaborator on server.server_id = server_collaborator.server_id").Where(" server.parent_id = ? AND (server_collaborator.user_id = ? OR server.create_by = ?)", s.ParentId, s.CreateBy, s.CreateBy)
	if s.ServerId != 0 {
		table = table.Where("server_id = ?", s.ServerId)
	}
	if s.Host != "" {
		table = table.Where("host = ?", s.Host)
	}
	if s.Name != "" {
		table = table.Where("name = ?", s.Name)
	}
	if err := table.Order("sort").Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&doc).Error; err != nil {
		return nil, 0, err
	}
	var count int
	table.Where("`deleted_at` IS NULL").Count(&count)
	return doc, count, nil
}

// Insert 添加服务器
func (s *ServerDao) Insert() (id int, err error) {
	s.UpdateBy = ""
	s.ServerId = 0
	result := model.DB.Table("server").Create(&s)
	if result.Error != nil {
		err = result.Error
		return
	}
	id = s.ServerId
	return
}

// BatchDelete 批量删除
func (s *ServerDao) BatchDelete(id []int) (err error) {
	if err = model.DB.Table("server").Where("server_id in (?)", id).Delete(&ServerDao{}).Error; err != nil {
		return
	}
	return
}

//Get 获取服务器信息
func (s *ServerDao) Get() (err error) {
	table := model.DB.Table("server")
	if s.ServerId != 0 {
		table = table.Where("server_id = ?", s.ServerId)
	}
	if err = table.First(&s).Error; err != nil {
		return
	}
	return
}

// 修改服务器信息
func (s *ServerDao) Update(serverId int) (update ServerDao, err error) {
	if err = model.DB.Table("server").First(&update, serverId).Error; err != nil {
		return
	}
	if err = model.DB.Table("server").Model(&update).Updates(&s).Error; err != nil {
		return
	}
	return
}

//校验是否有权限
func (s *ServerDao) CheckPermission(serverId, userId int) bool {
	var count int
	model.DB.Select("server.*").Table("server").Joins("LEFT JOIN server_collaborator on server.server_id = server_collaborator.server_id").Where("server.server_id = ?", serverId).Where(" server_collaborator.user_id = ? OR server.create_by = ?", userId, userId).Count(&count)
	if count == 0 {
		return false
	}
	return true
}
