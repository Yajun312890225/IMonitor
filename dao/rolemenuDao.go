package dao

import (
	"fmt"
	"iMonitor/model"
	"iMonitor/pkg/casbin"
)

// RolemenuDao 对role模型进行增删查改的单例工具
type RolemenuDao struct {
	model.RoleMenu
}

var rolemenuDao *RolemenuDao

// RoleMenu 得到dao-rolemenu 单例工具
func RoleMenu() *RolemenuDao {
	if rolemenuDao == nil {
		rolemenuDao = &RolemenuDao{}
	}
	return rolemenuDao
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
		cas.Enforce.AddPolicy(role.RoleKey, menu[i].Path, menu[i].Action)
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
	// sql3 := "delete from casbin_rule where v0= '" + role.RoleKey + "';"
	// orm.Eloquent.Exec(sql3)
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
