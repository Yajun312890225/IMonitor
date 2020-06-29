package dao

import (
	"errors"
	"iMonitor/model"
)

type ServerCollaboratorDao model.ServerCollaborator

type ReqCollaborator struct {
	ServerId int    `json:"serverId" binding:"required"` //当前服务器ID
	Username string `json:"username" binding:"required"` //协作者名字
}

func ServerCollaborator() *ServerCollaboratorDao {
	return &ServerCollaboratorDao{}
}

// AddCollaborator 添加协作者
func (sc *ServerCollaboratorDao) AddCollaborator() (err error) {
	var count int
	model.DB.Table("server_collaborator").Where("server_id = ? AND user_id = ? ", sc.ServerId, sc.UserId).Count(&count)
	if count != 0 {
		err = errors.New("协作者已经存在")
		return
	}

	err = model.DB.Table("server_collaborator").Create(&sc).Error
	return
}

// DelCollaborator 删除协作者
func (sc *ServerCollaboratorDao) DelCollaborator(serverId, userId int) (err error) {
	err = model.DB.Table("server_collaborator").Where("server_id = ? AND user_id = ? ", serverId, userId).Delete(&model.ServerCollaborator{}).Error
	return
}

// GetCollaborator 获取服务器的协作者
func (sc *ServerCollaboratorDao) GetCollaborator(serverId int) (users []UserDao, err error) {
	err = model.DB.Table("server_collaborator").Where("server_id = ?", serverId).Find(&users).Error
	return
}
