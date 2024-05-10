package service

import (
	"container/list"
	"errors"
	"github.com/MarchGe/go-admin-server/app/admin/model"
	"github.com/MarchGe/go-admin-server/app/admin/service/dto/req"
	"github.com/MarchGe/go-admin-server/app/admin/service/dto/res"
	"github.com/MarchGe/go-admin-server/app/common/E"
	"github.com/MarchGe/go-admin-server/app/common/database"
	"gorm.io/gorm"
	"sort"
	"time"
)

var _deptService = &DeptService{}

type DeptService struct {
}

func GetDeptService() *DeptService {
	return _deptService
}

func (s *DeptService) CreateDept(m *req.DeptUpsertReq) error {
	if m.ParentId > 0 {
		parentDept, _ := s.FindOneById(m.ParentId)
		if parentDept == nil {
			return E.Message("父级部门不存在")
		}
	}
	dept := s.toModel(m)
	dept.CreateTime = time.Now()
	dept.UpdateTime = time.Now()
	return database.GetMysql().Save(dept).Error
}

func (s *DeptService) toModel(req *req.DeptUpsertReq) *model.Dept {
	return &model.Dept{
		Name:     req.Name,
		SortNum:  req.SortNum,
		ParentId: req.ParentId,
	}
}

func (s *DeptService) FindOneById(id int64) (*model.Dept, error) {
	m := &model.Dept{}
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

func (s *DeptService) UpdateDept(id int64, m *req.DeptUpsertReq) error {
	if m.ParentId == id {
		return E.Message("父级部门有误")
	}
	if m.ParentId > 0 {
		parentDept, _ := s.FindOneById(m.ParentId)
		if parentDept == nil {
			return E.Message("父级部门不存在")
		}
	}
	dept, _ := s.FindOneById(id)
	if dept == nil {
		return E.Message("操作的部门不存在")
	}
	s.copyProperties(m, dept)
	dept.UpdateTime = time.Now()
	return database.GetMysql().Save(dept).Error
}

func (s *DeptService) copyProperties(req *req.DeptUpsertReq, dept *model.Dept) {
	dept.Name = req.Name
	dept.SortNum = req.SortNum
	dept.ParentId = req.ParentId
}

func (s *DeptService) DeleteDept(id int64) error {
	dept, _ := s.FindOneById(id)
	if dept == nil {
		return E.Message("操作的部门不存在")
	}
	err := database.GetMysql().Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.Dept{}, id).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *DeptService) FindDeptTree(keyword string) ([]*res.DeptTree, error) {
	depts := make([]*model.Dept, 0)
	db := database.GetMysql().Model(&model.Dept{})
	if keyword != "" {
		db.Where("name like ?", "%"+keyword+"%")
	}
	err := db.Find(&depts).Error
	if err != nil {
		return nil, err
	}
	return s.transferToDeptTree(depts), nil
}

func (s *DeptService) transferToDeptTree(depts []*model.Dept) []*res.DeptTree {
	indexMap := make(map[int64]*res.DeptTree)
	for _, dept := range depts {
		r := &res.DeptTree{}
		r.Dept = *dept
		r.Children = make([]*res.DeptTree, 0)
		indexMap[dept.Id] = r
	}
	var removeKeys = list.New()
	for key, dept := range indexMap {
		deptTree, ok := indexMap[dept.ParentId]
		if ok {
			deptTree.Children = append(deptTree.Children, dept)
			removeKeys.PushBack(key)
		}
	}
	for element := removeKeys.Front(); element != nil; element = element.Next() {
		key := element.Value.(int64)
		delete(indexMap, key)
	}

	resList := make([]*res.DeptTree, 0, len(indexMap))
	for _, v := range indexMap {
		resList = append(resList, v)
	}
	s.sortDeptTree(resList)
	return resList
}

func (s *DeptService) sortDeptTree(trees []*res.DeptTree) {
	if len(trees) == 0 {
		return
	}
	sort.Slice(trees, func(i, j int) bool {
		return trees[j].SortNum > trees[i].SortNum
	})
	for _, item := range trees {
		s.sortDeptTree(item.Children)
	}
}

func (s *DeptService) ExistsByParentId(parentId int64) (bool, error) {
	var existCount int64 = 0
	err := database.GetMysql().Model(&model.Dept{}).Select("id").Where("parent_id = ?", parentId).Limit(1).Count(&existCount).Error
	if err != nil {
		return false, err
	}
	return existCount > 0, nil
}
