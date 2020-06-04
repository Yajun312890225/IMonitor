package dao

import (
	"fmt"
	"iMonitor/model"
	"iMonitor/pkg/casbin"
)

// RolemenuDao 对role模型进行增删查改的
type RolemenuDao struct {
	model.RoleMenu
}

// RoleMenu
func RoleMenu() *RolemenuDao {
	return &RolemenuDao{}
}

// Insert 插入角色对应的菜单
func (rm *RolemenuDao) Insert(roleId int, menuId []int) error {
	var role model.Role
	if err := model.DB.Table("role").Where("role_id = ?", roleId).First(&role).Error; err != nil {
		return err
	}
	var menu []model.Menu
	if err := model.DB.Table("menu").Where("menu_id in (?)", menuId).Find(&menu).Error; err != nil {
		return err
	}
	sql := "INSERT INTO `role_menu` (`role_id`,`menu_id`,`role_name`) VALUES "
	cas := casbin.GetCasbin()
	for i := 0; i < len(menu); i++ {
		if len(menu)-1 == i {
			sql += fmt.Sprintf("(%d,%d,'%s');", role.RoleId, menu[i].MenuId, role.RoleKey)
		} else {
			sql += fmt.Sprintf("(%d,%d,'%s'),", role.RoleId, menu[i].MenuId, role.RoleKey)
		}
		if menu[i].MenuType == "A" {
			cas.Enforce.AddPolicy(role.RoleKey, menu[i].Path, menu[i].Action)
		}
	}
	model.DB.Exec(sql)
	return nil
}

// DeleteRoleMenu 删除角色对应的菜单
func (rm *RolemenuDao) DeleteRoleMenu(roleId int) error {

	if err := model.DB.Table("role_menu").Where("role_id = ?", roleId).Delete(&rm).Error; err != nil {
		return err
	}
	var role model.Role
	if err := model.DB.Table("role").Where("role_id = ?", roleId).First(&role).Error; err != nil {
		return err
	}

	var rules []model.CasbinRule
	if err := model.DB.Table("casbin_rule").Where("v0= ?", role.RoleKey).Find(&rules).Error; err != nil {
		return err
	}
	cas := casbin.GetCasbin()
	for _, rule := range rules {
		cas.Enforce.RemovePolicy(rule.V0, rule.V1, rule.V2)
	}
	return nil

}

// BatchDeleteRoleMenu 批量删除角色对应菜单
func (rm *RolemenuDao) BatchDeleteRoleMenu(roleIds []int) error {
	if err := model.DB.Table("role_menu").Where("role_id in (?)", roleIds).Delete(&rm).Error; err != nil {
		return err
	}
	var roles []model.Role
	if err := model.DB.Table("role").Where("role_id in (?)", roleIds).Find(&roles).Error; err != nil {
		return err
	}
	for _, role := range roles {
		var rules []model.CasbinRule
		if err := model.DB.Table("casbin_rule").Where("v0= ?", role.RoleKey).Find(&rules).Error; err != nil {
			return err
		}
		cas := casbin.GetCasbin()
		for _, rule := range rules {
			cas.Enforce.RemovePolicy(rule.V0, rule.V1, rule.V2)
		}
	}
	return nil
}

// GetPermis 获取权限
func (rm *RolemenuDao) GetPermis() ([]string, error) {
	var r []model.Menu
	table := model.DB.Select("menu.permission").Table("menu").Joins("left join role_menu on menu.menu_id = role_menu.menu_id")

	table = table.Where("role_id = ?", rm.RoleId)

	table = table.Where("menu.menu_type in('F')")
	if err := table.Find(&r).Error; err != nil {
		return nil, err
	}
	var list []string
	for i := 0; i < len(r); i++ {
		list = append(list, r[i].Permission)
	}
	return list, nil
}
