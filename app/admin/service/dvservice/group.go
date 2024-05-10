package dvservice

import (
	"errors"
	"fmt"
	"github.com/MarchGe/go-admin-server/app/admin/model/dvmodel"
	"github.com/MarchGe/go-admin-server/app/admin/service/dto/res"
	"github.com/MarchGe/go-admin-server/app/admin/service/dvservice/dto/req"
	"github.com/MarchGe/go-admin-server/app/common/E"
	"github.com/MarchGe/go-admin-server/app/common/database"
	"gorm.io/gorm"
	"time"
)

var _groupService = &GroupService{}

type GroupService struct {
}

func GetGroupService() *GroupService {
	return _groupService
}

func (s *GroupService) CreateGroup(info *req.GroupUpsertReq) error {
	group := s.toModel(info)
	group.UpdateTime = time.Now()
	group.CreateTime = time.Now()
	err := database.GetMysql().Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("HostList").Save(group).Error; err != nil {
			return err
		}
		if err := s.updateHostGroupRelations(tx, group.Id, info.HostIds, false); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *GroupService) toModel(info *req.GroupUpsertReq) *dvmodel.Group {
	return &dvmodel.Group{
		Name:    info.Name,
		SortNum: info.SortNum,
	}
}

func (s *GroupService) UpdateGroup(id int64, info *req.GroupUpsertReq) error {
	group, _ := s.FindOneById(id)
	if group == nil {
		return E.Message("操作的分组不存在")
	}
	s.copyProperties(info, group)
	group.UpdateTime = time.Now()
	err := database.GetMysql().Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("HostList").Save(group).Error; err != nil {
			return err
		}
		if err := s.updateHostGroupRelations(tx, group.Id, info.HostIds, true); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *GroupService) FindOneById(id int64, preloads ...string) (*dvmodel.Group, error) {
	m := &dvmodel.Group{}
	db := database.GetMysql()
	if len(preloads) > 0 {
		for _, preload := range preloads {
			db = db.Preload(preload)
		}
	}
	err := db.First(m, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, E.Message("查询的记录不存在")
		} else {
			return nil, err
		}
	}
	return m, nil
}

func (s *GroupService) copyProperties(info *req.GroupUpsertReq, group *dvmodel.Group) {
	group.Name = info.Name
	group.SortNum = info.SortNum
}

func (s *GroupService) DeleteGroup(id int64) error {
	err := database.GetMysql().Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&dvmodel.Group{}, id).Error; err != nil {
			return err
		}
		if err := tx.Where("group_id = ?", id).Delete(&dvmodel.HostGroup{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *GroupService) PageList(keyword string, page, pageSize int) (*res.PageableData[*dvmodel.Group], error) {
	groups := make([]*dvmodel.Group, 0)
	pageableData := &res.PageableData[*dvmodel.Group]{}
	db := database.GetMysql().Model(&dvmodel.Group{}).Preload("HostList")
	if keyword != "" {
		db.Where("name like ?", "%"+keyword+"%")
	}
	var count int64
	err := db.Count(&count).Order("sort_num").Offset(pageSize * (page - 1)).Limit(pageSize).Find(&groups).Error
	if err != nil {
		return nil, err
	}
	pageableData.List = groups
	pageableData.Total = count
	return pageableData, nil
}

func (s *GroupService) updateHostGroupRelations(tx *gorm.DB, groupId int64, hostIds []int64, update bool) error {
	if update {
		if err := tx.Where("group_id = ?", groupId).Delete(&dvmodel.HostGroup{}).Error; err != nil {
			return fmt.Errorf("delete host group relations error, %w", err)
		}
	}
	relations := make([]*dvmodel.HostGroup, len(hostIds))
	for i := range hostIds {
		relations[i] = &dvmodel.HostGroup{
			HostId:  hostIds[i],
			GroupId: groupId,
		}
	}
	if err := tx.Save(relations).Error; err != nil {
		return fmt.Errorf("save host group relations error, %w", err)
	}
	return nil
}
