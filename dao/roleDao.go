package dao

import "iMonitor/model"

// RoleDao 对role模型进行增删查改的单例工具
type RoleDao struct {
	model.Role
}

var roleDao *RoleDao

// Role 得到dao-role 单例工具
func Role() *RoleDao {
	if roleDao == nil {
		roleDao = &RoleDao{}
	}
	return roleDao
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
	if err := model.DB.Table("role_menu").Select("role_menu.menu_id").Joins("LEFT JOIN menu on menu.menu_id=role_menu.menu_id").Where("role_id = ? ", r.RoleId).Where(" role_menu.menu_id not in(select menu.parent_id from role_menu LEFT JOIN menu on menu.menu_id=role_menu.menu_id where role_id =? )", r.RoleId).Find(&menuList).Error; err != nil {
		return nil, err
	}
	for i := 0; i < len(menuList); i++ {
		menuIds = append(menuIds, menuList[i].MenuId)
	}
	return menuIds, nil
}
