package service

import (
	"errors"
	"fmt"
	"github.com/MarchGe/go-admin-server/app/admin/model"
	"github.com/MarchGe/go-admin-server/app/admin/service/dto/req"
	"github.com/MarchGe/go-admin-server/app/admin/service/dto/res"
	"github.com/MarchGe/go-admin-server/app/common/E"
	"github.com/MarchGe/go-admin-server/app/common/database"
	"github.com/MarchGe/go-admin-server/app/common/middleware/authz"
	"gorm.io/gorm"
	"time"
)

var _roleService = &RoleService{}

type RoleService struct {
}

func GetRoleService() *RoleService {
	return _roleService
}

func (s *RoleService) CreateRole(info *req.RoleUpsertReq) error {
	role := s.toModel(info)
	role.UpdateTime = time.Now()
	role.CreateTime = time.Now()
	return database.GetMysql().Omit("MenuList").Save(role).Error
}

func (s *RoleService) toModel(info *req.RoleUpsertReq) *model.Role {
	return &model.Role{
		Name:    info.Name,
		SortNum: info.SortNum,
	}
}

func (s *RoleService) UpdateRole(id int64, info *req.RoleUpsertReq) error {
	role, _ := s.FindOneById(id)
	if role == nil {
		return E.Message("操作的角色不存在")
	}
	s.copyProperties(info, role)
	role.UpdateTime = time.Now()
	return database.GetMysql().Omit("MenuList").Save(role).Error
}

func (s *RoleService) FindOneById(id int64) (*model.Role, error) {
	m := &model.Role{}
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

func (s *RoleService) copyProperties(info *req.RoleUpsertReq, role *model.Role) {
	role.Name = info.Name
	role.SortNum = info.SortNum
}

func (s *RoleService) DeleteRole(id int64) error {
	err := database.GetMysql().Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.Role{}, id).Error; err != nil {
			return err
		}
		if err := tx.Where("role_id = ?", id).Delete(&model.RoleUser{}).Error; err != nil {
			return err
		}
		if err := tx.Where("role_id = ?", id).Delete(&model.RoleMenu{}).Error; err != nil {
			return err
		}
		return authz.DeleteRole(authz.RoleSub(id))
	})
	return err
}

func (s *RoleService) PageList(keyword string, page, pageSize int) (*res.PageableData[*model.Role], error) {
	roles := make([]*model.Role, 0)
	pageableData := &res.PageableData[*model.Role]{}
	db := database.GetMysql().Model(&model.Role{})
	if keyword != "" {
		db.Where("name like ?", "%"+keyword+"%")
	}
	var count int64
	err := db.Count(&count).Order("sort_num").Offset(pageSize * (page - 1)).Limit(pageSize).Preload("MenuList").Find(&roles).Error
	if err != nil {
		return nil, err
	}
	pageableData.List = roles
	pageableData.Total = count
	return pageableData, nil
}

func (s *RoleService) UpdateRoleMenus(id int64, updateReq *req.RoleMenusUpdateReq) error {
	err := database.GetMysql().Transaction(func(tx *gorm.DB) error {
		// 删除role和menu的关联关系
		if err := tx.Where("role_id = ?", id).Delete(&model.RoleMenu{}).Error; err != nil {
			return err
		}
		if len(updateReq.Ids) > 0 {
			// 添加role和menu的关联关系
			menus := make([]*model.RoleMenu, 0, len(updateReq.Ids))
			for _, menuId := range updateReq.Ids {
				menus = append(menus, &model.RoleMenu{
					RoleId: id,
					MenuId: menuId,
				})
			}
			if err := tx.Save(menus).Error; err != nil {
				return err
			}
		}
		symbols, err := GetMenuService().FindSymbolsByIds(updateReq.Ids)
		if err != nil {
			return fmt.Errorf("find menu symbols error, %w", err)
		}
		return authz.UpdateSubPolicies(authz.RoleSub(id), symbols)
	})
	return err
}
