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

var _scriptService = &ScriptService{}

type ScriptService struct {
}

func GetScriptService() *ScriptService {
	return _scriptService
}

func (s *ScriptService) CreateScript(info *req.ScriptUpsertReq) error {
	existScript, _ := s.FindOneByNameAndVersion(info.Name, info.Version)
	if existScript != nil {
		return E.Message(fmt.Sprintf("脚本'%s:%s'已存在", info.Name, info.Version))
	}
	script := s.toModel(info)
	script.CreateTime = time.Now()
	script.UpdateTime = time.Now()
	err := database.GetMysql().Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(script).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *ScriptService) toModel(info *req.ScriptUpsertReq) *dvmodel.Script {
	return &dvmodel.Script{
		Name:        info.Name,
		Version:     info.Version,
		Content:     info.Content,
		Description: info.Description,
	}
}

func (s *ScriptService) UpdateScript(id int64, info *req.ScriptUpsertReq) error {
	script, _ := s.FindOneById(id)
	if script == nil {
		return E.Message("操作的脚本不存在")
	}
	existScript, _ := s.FindOneByNameAndVersion(info.Name, info.Version)
	if existScript != nil && existScript.Id != id {
		return E.Message(fmt.Sprintf("脚本'%s:%s'已存在", info.Name, info.Version))
	}
	s.copyProperties(info, script)
	script.UpdateTime = time.Now()
	err := database.GetMysql().Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(script).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *ScriptService) FindOneById(id int64) (*dvmodel.Script, error) {
	m := &dvmodel.Script{}
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

func (s *ScriptService) copyProperties(info *req.ScriptUpsertReq, script *dvmodel.Script) {
	script.Name = info.Name
	script.Version = info.Version
	script.Content = info.Content
	script.Description = info.Description
}

func (s *ScriptService) DeleteScript(id int64) error {
	script, _ := s.FindOneById(id)
	if script == nil {
		return E.Message("操作的脚本不存在")
	}
	err := database.GetMysql().Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&dvmodel.Script{}, id).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *ScriptService) PageList(keyword string, page, pageSize int) (*res.PageableData[*dvmodel.Script], error) {
	scripts := make([]*dvmodel.Script, 0)
	pageableData := &res.PageableData[*dvmodel.Script]{}
	db := database.GetMysql().Model(&dvmodel.Script{})
	if keyword != "" {
		db.Where("name like ?", "%"+keyword+"%")
	}
	var count int64
	err := db.Count(&count).Order("create_time desc").Offset(pageSize * (page - 1)).Limit(pageSize).Find(&scripts).Error
	if err != nil {
		return nil, err
	}
	pageableData.List = scripts
	pageableData.Total = count
	return pageableData, nil
}

func (s *ScriptService) FindOneByNameAndVersion(name, version string) (*dvmodel.Script, error) {
	m := &dvmodel.Script{}
	err := database.GetMysql().Where("name = ? and version = ?", name, version).First(m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, E.Message("查询的记录不存在")
		} else {
			return nil, err
		}
	}
	return m, nil
}
