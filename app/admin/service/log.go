package service

import (
	"github.com/MarchGe/go-admin-server/app/admin/model"
	"github.com/MarchGe/go-admin-server/app/admin/service/dto/res"
	"github.com/MarchGe/go-admin-server/app/common/database"
	"time"
)

var _logService = &LogService{}

type LogService struct {
}

func GetLogService() *LogService {
	return _logService
}

func (s *LogService) AddLoginLog(log *model.LoginLog) error {
	if err := database.GetMysql().Save(log).Error; err != nil {
		return err
	}
	return nil
}
func (s *LogService) AddOpLog(log *model.OpLog) error {
	if err := database.GetMysql().Save(log).Error; err != nil {
		return err
	}
	return nil
}
func (s *LogService) AddExceptionLog(log *model.ExceptionLog) error {
	if err := database.GetMysql().Save(log).Error; err != nil {
		return err
	}
	return nil
}

func (s *LogService) LoginLogPageList(keyword string, page, pageSize int) (*res.PageableData[*model.LoginLog], error) {
	logs := make([]*model.LoginLog, 0)
	pageableData := &res.PageableData[*model.LoginLog]{}
	db := database.GetMysql().Model(&model.LoginLog{})
	if keyword != "" {
		db.Where("nickname like ?", "%"+keyword+"%")
	}
	var count int64
	err := db.Count(&count).Order("create_time desc").Offset(pageSize * (page - 1)).Limit(pageSize).Find(&logs).Error
	if err != nil {
		return nil, err
	}
	pageableData.List = logs
	pageableData.Total = count
	return pageableData, nil
}
func (s *LogService) OpLogPageList(keyword string, page, pageSize int) (*res.PageableData[*model.OpLog], error) {
	logs := make([]*model.OpLog, 0)
	pageableData := &res.PageableData[*model.OpLog]{}
	db := database.GetMysql().Model(&model.OpLog{})
	if keyword != "" {
		db.Where("nickname like ?", "%"+keyword+"%")
	}
	var count int64
	err := db.Count(&count).Order("create_time desc").Offset(pageSize * (page - 1)).Limit(pageSize).Find(&logs).Error
	if err != nil {
		return nil, err
	}
	pageableData.List = logs
	pageableData.Total = count
	return pageableData, nil
}
func (s *LogService) ExceptionLogPageList(keyword string, page, pageSize int) (*res.PageableData[*model.ExceptionLog], error) {
	logs := make([]*model.ExceptionLog, 0)
	pageableData := &res.PageableData[*model.ExceptionLog]{}
	db := database.GetMysql().Model(&model.ExceptionLog{})
	if keyword != "" {
		db.Where("nickname like ?", "%"+keyword+"%")
	}
	var count int64
	err := db.Count(&count).Order("create_time desc").Offset(pageSize * (page - 1)).Limit(pageSize).Find(&logs).Error
	if err != nil {
		return nil, err
	}
	pageableData.List = logs
	pageableData.Total = count
	return pageableData, nil
}

func (s *LogService) DeleteLoginLog(beforeTime time.Time) error {
	for {
		result := database.GetMysql().Where("create_time <= ?", beforeTime).Limit(100).Delete(&model.LoginLog{})
		if err := result.Error; err != nil {
			return err
		}
		if result.RowsAffected == 0 {
			break
		}
	}
	return nil
}

func (s *LogService) DeleteOpLog(beforeTime time.Time) error {
	for {
		result := database.GetMysql().Where("create_time <= ?", beforeTime).Limit(100).Delete(&model.OpLog{})
		if err := result.Error; err != nil {
			return err
		}
		if result.RowsAffected == 0 {
			break
		}
	}
	return nil
}

func (s *LogService) DeleteExceptionLog(beforeTime time.Time) error {
	for {
		result := database.GetMysql().Where("create_time <= ?", beforeTime).Limit(100).Delete(&model.ExceptionLog{})
		if err := result.Error; err != nil {
			return err
		}
		if result.RowsAffected == 0 {
			break
		}
	}
	return nil
}
