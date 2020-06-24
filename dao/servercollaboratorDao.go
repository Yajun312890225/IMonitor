package dao

import "iMonitor/model"

type ServerCollaboratorDao model.ServerCollaborator

type ReqCreateCollaborator struct {
	ServerId int    `json:"serverId" binding:"required"` //当前服务器ID
	Username string `json:"username" binding:"required"` //协作者名字
}

func ServerCollaborator() *ServerCollaboratorDao {
	return &ServerCollaboratorDao{}
}

// AddCollaborator 添加协作者
func (sc *ServerCollaboratorDao) AddCollaborator() (err error) {
	err = model.DB.Table("server_collaborator").Create(&sc).Error
	return
}

// DelCollaborator 删除协作者
func (sc *ServerCollaboratorDao) DelCollaborator(serverId, userId int) (err error) {
	err = model.DB.Table("server_collaborator").Where("server_id = ? AND user_id = ? ", serverId, userId).Delete(&model.ServerCollaborator{}).Error
	return
}
