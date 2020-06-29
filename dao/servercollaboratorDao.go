package dao

import (
	"errors"
	"iMonitor/model"
	"strconv"
)

type ServerCollaboratorDao model.ServerCollaborator

type ReqCollaborator struct {
	ServerId int      `json:"serverId" binding:"required"` //当前服务器ID
	Username []string `json:"username" binding:"required"` //协作者名字
}

func ServerCollaborator() *ServerCollaboratorDao {
	return &ServerCollaboratorDao{}
}

// AddCollaborator 添加协作者
func (sc *ServerCollaboratorDao) AddCollaborator(serverId int, userId []int) (err error) {
	if len(userId) == 0 {
		err = errors.New("协作者不存在")
		return
	}
	for _, v := range userId {
		var count int
		model.DB.Table("server_collaborator").Where("server_id = ? AND user_id = ? ", serverId, v).Count(&count)
		if count != 0 {
			err = errors.New("部分协作者已经存在")
			return
		}
		sc.ServerId = serverId
		sc.UserId = v
		if err = model.DB.Table("server_collaborator").Create(&sc).Error; err != nil {
			return
		}
	}
	return
}

// DelCollaborator 删除协作者
func (sc *ServerCollaboratorDao) DelCollaborator(serverId int, userId []int) (err error) {
	if len(userId) == 0 {
		err = errors.New("协作者不存在")
		return
	}
	for _, v := range userId {
		var count int
		model.DB.Table("server_collaborator").Where("server_id = ? AND user_id = ? ", serverId, v).Count(&count)
		if count == 0 {
			err = errors.New("部分协作者不存在")
			return
		}
		if err = model.DB.Table("server_collaborator").Where("server_id = ? AND user_id = ? ", serverId, v).Delete(&model.ServerCollaborator{}).Error; err != nil {
			return
		}
	}
	return
}

// GetCollaborator 获取服务器的协作者
func (sc *ServerCollaboratorDao) GetCollaborator(serverId int) (users []UserDao, err error) {
	data := []ServerCollaboratorDao{}
	if err = model.DB.Table("server").Select("server_collaborator.*").Joins("left join server_collaborator on server.server_id=server_collaborator.server_id").Where("server.server_id = ?", serverId).Find(&data).Error; err != nil {
		return
	}
	userIds := func(sc []ServerCollaboratorDao) []int {
		userId := make([]int, 0)
		for _, v := range sc {
			userId = append(userId, v.UserId)
		}
		return userId
	}(data)
	var server ServerDao
	if err = model.DB.Table("server").Where("server_id = ?", serverId).First(&server).Error; err != nil {
		return
	}
	createBy, _ := strconv.Atoi(server.CreateBy)
	userIds = append(userIds, createBy)
	if err = model.DB.Table("user").Where("user_id in (?)", userIds).Find(&users).Error; err != nil {
		return
	}
	return
}
