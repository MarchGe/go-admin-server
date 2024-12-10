package dvservice

import (
	"errors"
	"fmt"
	"github.com/MarchGe/go-admin-server/app/admin/model/dvmodel"
	"github.com/MarchGe/go-admin-server/app/admin/service/dto/res"
	"github.com/MarchGe/go-admin-server/app/admin/service/dvservice/dto/req"
	dvRes "github.com/MarchGe/go-admin-server/app/admin/service/dvservice/dto/res"
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
	count, err := s.getRefCount(id)
	if err != nil {
		return err
	}
	if count > 0 {
		return E.Message("该脚本正在被引用，无法删除")
	}
	err = database.GetMysql().Transaction(func(tx *gorm.DB) error {
		if e := tx.Delete(&dvmodel.Script{}, id).Error; e != nil {
			return e
		}
		return nil
	})
	return err
}

func (s *ScriptService) PageList(keyword string, page, pageSize int) (*res.PageableData[*dvRes.ScriptRes], error) {
	scripts := make([]*dvmodel.Script, 0)
	db := database.GetMysql().Model(&dvmodel.Script{})
	if keyword != "" {
		db.Where("name like ?", "%"+keyword+"%")
	}
	var count int64
	err := db.Count(&count).Order("create_time desc").Offset(pageSize * (page - 1)).Limit(pageSize).Find(&scripts).Error
	if err != nil {
		return nil, err
	}
	ids := make([]int64, len(scripts))
	for i := range scripts {
		ids[i] = scripts[i].Id
	}
	scriptTaskRefMap, err := s.getScriptTaskRefMap(ids)
	deployTaskRefMap, err := s.getDeployTaskRefMap(ids)
	scriptResList := make([]*dvRes.ScriptRes, len(scripts))
	for i := range scripts {
		scriptResList[i] = &dvRes.ScriptRes{
			Script: *scripts[i],
		}
		cnt := scriptTaskRefMap[scripts[i].Id]
		scriptResList[i].ScriptTaskRefCount = cnt
		cnt = deployTaskRefMap[scripts[i].Id]
		scriptResList[i].DeployTaskRefCount = cnt
	}
	pageableData := &res.PageableData[*dvRes.ScriptRes]{}
	pageableData.List = scriptResList
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

func (s *ScriptService) getScriptTaskRefMap(ids []int64) (map[int64]int32, error) {
	results := make([]struct {
		Id       int64
		RefCount int32
	}, 0)
	err := database.GetMysql().Table("dv_script s").Select("s.id, count(s.id) as ref_count").Where("s.id in ?", ids).
		Joins("inner join dv_script_task_script sts on sts.script_id = s.id").
		Group("s.id").Scan(&results).Error
	if err != nil {
		return nil, err
	}
	refMap := make(map[int64]int32, len(results))
	for _, result := range results {
		refMap[result.Id] = result.RefCount
	}
	return refMap, nil
}

func (s *ScriptService) getDeployTaskRefMap(ids []int64) (map[int64]int32, error) {
	results := make([]struct {
		Id       int64
		RefCount int32
	}, 0)
	err := database.GetMysql().Table("dv_script s").Select("s.id, count(s.id) as ref_count").Where("s.id in ?", ids).
		Joins("inner join dv_task_deploy td on td.script_id = s.id").
		Group("s.id").Scan(&results).Error
	if err != nil {
		return nil, err
	}
	refMap := make(map[int64]int32, len(results))
	for _, result := range results {
		refMap[result.Id] = result.RefCount
	}
	return refMap, nil
}

func (s *ScriptService) getRefCount(id int64) (int32, error) {
	var ids = []int64{id}
	dMap, e1 := s.getDeployTaskRefMap(ids)
	sMap, e2 := s.getScriptTaskRefMap(ids)
	if e := errors.Join(e1, e2); e != nil {
		return 0, fmt.Errorf("get ref count error, %w", e)
	}
	dCnt := dMap[id]
	sCnt := sMap[id]
	return dCnt + sCnt, nil
}
