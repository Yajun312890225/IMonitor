package dao

import "iMonitor/model"

type ServerDao model.Server
type ReqPing struct {
	Host string `json:"host" binding:"required"`
	Port string `json:"port" binding:"required"`
	Key1 string `json:"key1" binding:"required"`
	Key2 string `json:"key2" binding:"required"`
}

func Server() *ServerDao {
	return &ServerDao{}
}

//GetPage 获取服务器
func (s *ServerDao) GetPage(pageSize int, pageIndex int) ([]model.Server, int, error) {
	var doc []model.Server

	table := model.DB.Select("*").Table("server")
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
