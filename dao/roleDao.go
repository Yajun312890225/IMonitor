package dao

import (
	"errors"
	"iMonitor/model"
)

// RoleDao 对role模型进行增删查改
type RoleDao struct {
	model.Role
}

// Role
func Role() *RoleDao {
	return &RoleDao{}
}

type MenuIdList struct {
	MenuId int `json:"menuId"`
}

//GetPage 获取角色页
func (r *RoleDao) GetPage(pageSize int, pageIndex int) ([]model.Role, int, error) {
	var doc []model.Role

	table := model.DB.Select("*").Table("role")
	if r.RoleId != 0 {
		table = table.Where("role_id = ?", r.RoleId)
	}
	if r.RoleName != "" {
		table = table.Where("role_name = ?", r.RoleName)
	}
	if r.Status != "" {
		table = table.Where("status = ?", r.Status)
	}
	if r.RoleKey != "" {
		table = table.Where("role_key = ?", r.RoleKey)
	}
	if err := table.Order("role_sort").Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&doc).Error; err != nil {
		return nil, 0, err
	}
	var count int
	table.Where("`deleted_at` IS NULL").Count(&count)
	return doc, count, nil
}

//Get 获取角色
func (r *RoleDao) Get() (err error) {
	table := model.DB.Table("role")
	if r.RoleId != 0 {
		table = table.Where("role_id = ?", r.RoleId)
	}
	if r.RoleName != "" {
		table = table.Where("role_name = ?", r.RoleName)
	}
	if err = table.First(&r).Error; err != nil {
		return
	}
	return
}

// GetRoleMeunId 获取角色对应的菜单ids
func (r *RoleDao) GetRoleMeunId() ([]int, error) {
	menuIds := make([]int, 0)
	menuList := make([]MenuIdList, 0)
	// if err := model.DB.Table("role_menu").Select("role_menu.menu_id").Joins("LEFT JOIN menu on menu.menu_id=role_menu.menu_id").Where("role_id = ? ", r.RoleId).Where(" role_menu.menu_id not in(select menu.parent_id from role_menu LEFT JOIN menu on menu.menu_id=role_menu.menu_id where role_id =? )", r.RoleId).Find(&menuList).Error; err != nil {
	// 	return nil, err
	// }
	if err := model.DB.Table("role_menu").Select("role_menu.menu_id").Joins("LEFT JOIN menu on menu.menu_id=role_menu.menu_id").Where("role_id = ? ", r.RoleId).Find(&menuList).Error; err != nil {
		return nil, err
	}
	for i := 0; i < len(menuList); i++ {
		menuIds = append(menuIds, menuList[i].MenuId)
	}
	return menuIds, nil
}

// Insert 创建角色
func (r *RoleDao) Insert() (id int, err error) {
	// i := 0
	// model.DB.Table("role").Where("role_name=? or role_key = ?", r.RoleName, r.RoleKey).Count(&i)
	// if i > 0 {
	// 	return 0, errors.New("角色名称或者角色标识已经存在！")
	// }
	r.UpdateBy = ""
	r.RoleId = 0
	result := model.DB.Table("role").Create(&r)
	if result.Error != nil {
		err = result.Error
		return
	}
	id = r.RoleId
	return
}

// Update 修改角色
func (r *RoleDao) Update(id int) (update RoleDao, err error) {
	if err = model.DB.Table("role").First(&update, id).Error; err != nil {
		return
	}

	if r.RoleName != "" && r.RoleName != update.RoleName {
		return update, errors.New("角色名称不允许修改！")
	}

	if r.RoleKey != "" && r.RoleKey != update.RoleKey {
		return update, errors.New("角色标识不允许修改！")
	}

	if err = model.DB.Table("role").Model(&update).Updates(&r).Error; err != nil {
		return
	}
	return
}

// BatchDelete 批量删除
func (r *RoleDao) BatchDelete(id []int) (err error) {
	if err = model.DB.Table("role").Where("role_id in (?)", id).Delete(&RoleDao{}).Error; err != nil {
		return
	}
	return
}

func (r *RoleDao) GetList() (role []RoleDao, err error) {
	table := model.DB.Table("role")
	if r.RoleId != 0 {
		table = table.Where("role_id = ?", r.RoleId)
	}
	if r.RoleName != "" {
		table = table.Where("role_name = ?", r.RoleName)
	}
	if err = table.Order("role_sort").Find(&role).Error; err != nil {
		return
	}
	return
}
