package dao

import (
	"iMonitor/model"
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
