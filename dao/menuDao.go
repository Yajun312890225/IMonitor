package dao

import (
	"errors"
	"iMonitor/model"
	"strconv"
)

// MenuDao 对menu模型进行增删查改的单例工具
type MenuDao struct {
	model.Menu
}

var menuDao *MenuDao

// Menu 得到dao-menu 单例工具
func Menu() *MenuDao {
	if menuDao == nil {
		menuDao = &MenuDao{}
	}
	return menuDao
}

// GetAllMenu 获取所有Menu
func (m *MenuDao) GetAllMenu() (menu []MenuDao, err error) {
	menulist, err := m.GetPage()

	menu = make([]MenuDao, 0)
	for i := 0; i < len(menulist); i++ {
		if menulist[i].ParentId != 0 {
			continue
		}
		menusInfo := RecursionMenu(&menulist, menulist[i].Menu)

		menu = append(menu, MenuDao{menusInfo})
	}
	return
}

// RecursionMenu 递归查找Menu关系
func RecursionMenu(menulist *[]MenuDao, menu model.Menu) model.Menu {
	list := *menulist

	min := make([]model.Menu, 0)
	for j := 0; j < len(list); j++ {

		if menu.MenuId != list[j].ParentId {
			continue
		}
		mi := model.Menu{}
		mi.MenuId = list[j].MenuId
		mi.MenuName = list[j].MenuName
		mi.Title = list[j].Title
		mi.Icon = list[j].Icon
		mi.Path = list[j].Path
		mi.MenuType = list[j].MenuType
		mi.Action = list[j].Action
		mi.Permission = list[j].Permission
		mi.ParentId = list[j].ParentId
		mi.NoCache = list[j].NoCache
		mi.Breadcrumb = list[j].Breadcrumb
		mi.Component = list[j].Component
		mi.Sort = list[j].Sort
		mi.Visible = list[j].Visible
		mi.Children = []model.Menu{}

		if mi.MenuType != "F" {
			ms := RecursionMenu(menulist, mi)
			min = append(min, ms)

		} else {
			min = append(min, mi)
		}

	}
	menu.Children = min
	return menu
}

//GetPage 查找所有Menu信息等待去处理
func (m *MenuDao) GetPage() (menus []MenuDao, err error) {
	table := model.DB.Table("menu")
	if m.MenuName != "" {
		table = table.Where("menu_name = ?", m.MenuName)
	}
	if m.Title != "" {
		table = table.Where("title = ?", m.Title)
	}
	if m.Visible != "" {
		table = table.Where("visible = ?", m.Visible)
	}
	if m.MenuType != "" {
		table = table.Where("menu_type = ?", m.MenuType)
	}

	if err = table.Order("sort").Find(&menus).Error; err != nil {
		return
	}
	return
}

// Create 创建菜单
func (m *MenuDao) Create() (id int, err error) {
	m.MenuId = 0
	result := model.DB.Table("menu").Create(&m)
	if result.Error != nil {
		err = result.Error
		return
	}
	err = InitPaths(m)
	if err != nil {
		return
	}
	id = m.MenuId
	return
}

// InitPaths 初始化路径
func InitPaths(menu *MenuDao) (err error) {
	parentMenu := MenuDao{}
	if int(menu.ParentId) != 0 {
		model.DB.Table("menu").Where("menu_id = ?", menu.ParentId).First(&parentMenu)
		if parentMenu.Paths == "" {
			err = errors.New("父级paths异常，请尝试对当前节点父级菜单进行更新操作！")
			return
		}
		menu.Paths = parentMenu.Paths + "/" + strconv.Itoa(menu.MenuId)
	} else {
		menu.Paths = "/0/" + strconv.Itoa(menu.MenuId)
	}
	model.DB.Table("menu").Where("menu_id = ?", menu.MenuId).Update("paths", menu.Paths)
	return
}

// Update 更新菜单
func (m *MenuDao) Update(id int) (update MenuDao, err error) {
	if err = model.DB.Table("menu").First(&update, id).Error; err != nil {
		return
	}
	if err = model.DB.Table("menu").Model(&update).Updates(&m).Error; err != nil {
		return
	}
	err = InitPaths(m)
	if err != nil {
		return
	}
	return
}

// Delete 删除菜单
func (m *MenuDao) Delete(id int) (success bool, err error) {
	if err = model.DB.Table("menu").Where("menu_id = ?", id).Delete(&MenuDao{}).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}
