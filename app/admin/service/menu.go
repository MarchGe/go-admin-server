package service

import (
	"container/list"
	"errors"
	"fmt"
	"github.com/MarchGe/go-admin-server/app/admin/model"
	"github.com/MarchGe/go-admin-server/app/admin/service/dto/req"
	"github.com/MarchGe/go-admin-server/app/admin/service/dto/res"
	"github.com/MarchGe/go-admin-server/app/common/E"
	"github.com/MarchGe/go-admin-server/app/common/database"
	"github.com/MarchGe/go-admin-server/app/common/middleware/authz"
	"gorm.io/gorm"
	"slices"
	"sort"
	"strings"
	"time"
)

var _menuService = &MenuService{}

type MenuService struct {
}

func GetMenuService() *MenuService {
	return _menuService
}

func (s *MenuService) CreateMenu(m *req.MenuUpsertReq) error {
	if m.ParentId > 0 {
		parentMenu, _ := s.FindOneById(m.ParentId)
		if parentMenu == nil {
			return E.Message("父菜单不存在")
		}
	}
	if m.Symbol != "" {
		existMenu, _ := s.FindOneBySymbol(m.Symbol)
		if existMenu != nil {
			return E.Message(fmt.Sprintf("权限标识'%s'已存在", m.Symbol))
		}
	}
	menu := s.toModel(m)
	menu.CreateTime = time.Now()
	menu.UpdateTime = time.Now()
	return database.GetMysql().Save(menu).Error
}

func (s *MenuService) toModel(req *req.MenuUpsertReq) *model.Menu {
	return &model.Menu{
		Name:        req.Name,
		Icon:        req.Icon,
		SortNum:     req.SortNum,
		Url:         req.Url,
		Symbol:      req.Symbol,
		Display:     req.Display,
		External:    req.External,
		ExternalWay: req.ExternalWay,
		ParentId:    req.ParentId,
	}
}

func (s *MenuService) FindOneById(id int64) (*model.Menu, error) {
	m := &model.Menu{}
	err := database.GetMysql().First(m, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, E.Message("查询的记录不存在")
		} else {
			return nil, err
		}
	}
	return m, nil
}

func (s *MenuService) UpdateMenu(id int64, m *req.MenuUpsertReq) error {
	if m.ParentId == id {
		return E.Message("父菜单有误")
	}
	if m.ParentId > 0 {
		parentMenu, _ := s.FindOneById(m.ParentId)
		if parentMenu == nil {
			return E.Message("父菜单不存在")
		}
	}
	menu, _ := s.FindOneById(id)
	if menu == nil {
		return E.Message("操作的菜单不存在")
	}
	if m.Symbol != "" {
		existMenu, _ := s.FindOneBySymbol(m.Symbol)
		if existMenu != nil && existMenu.Id != id {
			return E.Message(fmt.Sprintf("权限标识'%s'已存在", m.Symbol))
		}
	}
	s.copyProperties(m, menu)
	menu.UpdateTime = time.Now()
	return database.GetMysql().Save(menu).Error
}

func (s *MenuService) FindOneBySymbol(symbol string) (*model.Menu, error) {
	m := &model.Menu{}
	err := database.GetMysql().Where("symbol = ?", symbol).First(m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, E.Message("查询的记录不存在")
		} else {
			return nil, err
		}
	}
	return m, nil
}

func (s *MenuService) copyProperties(req *req.MenuUpsertReq, menu *model.Menu) {
	menu.Name = req.Name
	menu.Url = req.Url
	menu.Icon = req.Icon
	menu.SortNum = req.SortNum
	menu.Display = req.Display
	menu.External = req.External
	menu.ExternalWay = req.ExternalWay
	menu.ParentId = req.ParentId
	menu.Symbol = req.Symbol
}

func (s *MenuService) DeleteMenu(id int64) error {
	menu, _ := s.FindOneById(id)
	if menu == nil {
		return E.Message("操作的菜单不存在")
	}
	err := database.GetMysql().Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.Menu{}, id).Error; err != nil {
			return err
		}
		if err := tx.Where("menu_id = ?", id).Delete(&model.RoleMenu{}).Error; err != nil {
			return err
		}
		if err := tx.Where("menu_id = ?", id).Delete(&model.UserMenu{}).Error; err != nil {
			return err
		}
		if menu.Symbol == "" {
			return nil
		}
		return authz.DeletePermission(menu.Symbol)
	})
	return err
}

func (s *MenuService) FindMenuTree() ([]*res.MenuTree, error) {
	menus := make([]*model.Menu, 0)
	err := database.GetMysql().Find(&menus).Error
	if err != nil {
		return nil, err
	}
	return s.transferToMenuTree(menus), nil
}

func (s *MenuService) FindUserMenuTree(userId int64, isRoot bool) ([]*res.MenuTree, error) {
	menus := make([]*model.Menu, 0)
	db := database.GetMysql()
	var err error
	if isRoot {
		err = db.Model(&model.Menu{}).Where("display = ?", 1).Find(&menus).Error
	} else {
		err = db.Model(&model.Menu{}).Where("id in (?)",
			db.Raw("? union ?",
				db.Model(&model.RoleMenu{}).Select("ops_menu.id").Joins("inner join ops_menu on ops_menu.id = ops_role_menu.menu_id").Where("ops_role_menu.role_id in (?) and ops_menu.display = 1",
					db.Model(&model.RoleUser{}).Select("role_id").Where("user_id = ?", userId),
				),
				db.Model(&model.UserMenu{}).Select("ops_menu.id").Joins("inner join ops_menu on ops_menu.id = ops_user_menu.menu_id").Where("ops_user_menu.user_id = ? and ops_menu.display = 1", userId),
			),
		).Find(&menus).Error
	}
	if err != nil {
		return nil, err
	}
	return s.transferToMenuTree(menus), nil
}

func (s *MenuService) transferToMenuTree(menus []*model.Menu) []*res.MenuTree {
	indexMap := make(map[int64]*res.MenuTree)
	for _, menu := range menus {
		r := &res.MenuTree{}
		r.Menu = *menu
		r.Children = make([]*res.MenuTree, 0)
		indexMap[menu.Id] = r
	}
	var removeKeys = list.New()
	for key, menu := range indexMap {
		menuTree, ok := indexMap[menu.ParentId]
		if ok {
			menuTree.Children = append(menuTree.Children, menu)
			removeKeys.PushBack(key)
		}
	}
	for element := removeKeys.Front(); element != nil; element = element.Next() {
		key := element.Value.(int64)
		delete(indexMap, key)
	}

	resList := make([]*res.MenuTree, 0, len(indexMap))
	for _, v := range indexMap {
		resList = append(resList, v)
	}
	s.sortMenuTree(resList)
	return resList
}

func (s *MenuService) sortMenuTree(trees []*res.MenuTree) {
	if len(trees) == 0 {
		return
	}
	sort.Slice(trees, func(i, j int) bool {
		return trees[j].SortNum > trees[i].SortNum
	})
	for _, item := range trees {
		s.sortMenuTree(item.Children)
	}
}

func (s *MenuService) FindSymbolsByIds(ids []int64) ([]string, error) {
	symbols := make([]string, 0)
	err := database.GetMysql().Model(&model.Menu{}).Select("symbol").Where("id in ?", ids).Find(&symbols).Error
	if err != nil {
		return nil, err
	}
	symbols = slices.DeleteFunc(symbols, func(s string) bool {
		return strings.TrimSpace(s) == ""
	})
	return symbols, nil
}

func (s *MenuService) ExistsByParentId(parentId int64) (bool, error) {
	var existCount int64 = 0
	err := database.GetMysql().Model(&model.Menu{}).Select("id").Where("parent_id = ?", parentId).Limit(1).Count(&existCount).Error
	if err != nil {
		return false, err
	}
	return existCount > 0, nil
}

func (s *MenuService) FindUserPermissions(userId int64, isRoot bool) ([]string, error) {
	permissions := make([]string, 0)
	db := database.GetMysql()
	var err error
	if isRoot {
		err = db.Model(&model.Menu{}).Select("Symbol").Scan(&permissions).Error
	} else {
		err = db.Model(&model.Menu{}).Select("Symbol").Where("id in (?)",
			db.Raw("? union ?",
				db.Model(&model.RoleMenu{}).Select("ops_menu.id").Joins("inner join ops_menu on ops_menu.id = ops_role_menu.menu_id").Where("ops_role_menu.role_id in (?)",
					db.Model(&model.RoleUser{}).Select("role_id").Where("user_id = ?", userId),
				),
				db.Model(&model.UserMenu{}).Select("ops_menu.id").Joins("inner join ops_menu on ops_menu.id = ops_user_menu.menu_id").Where("ops_user_menu.user_id = ?", userId),
			),
		).Scan(&permissions).Error
	}
	if err != nil {
		return nil, err
	}
	if len(permissions) == 0 {
		return permissions, nil
	}
	trimPermissions := make([]string, 0)
	for _, p := range permissions {
		if strings.TrimSpace(p) != "" {
			trimPermissions = append(trimPermissions, p)
		}
	}
	return trimPermissions, nil
}
