package service

import (
	"errors"
	"github.com/MarchGe/go-admin-server/app/admin/model"
	"github.com/MarchGe/go-admin-server/app/admin/service/dto/req"
	"github.com/MarchGe/go-admin-server/app/common/E"
	"github.com/MarchGe/go-admin-server/app/common/database"
	"gorm.io/gorm"
)

var _settingsService = &SettingsService{}

type SettingsService struct {
}

func GetSettingsService() *SettingsService {
	return _settingsService
}

func (s *SettingsService) UpsertSettings(userId int64, settings *req.SettingsUpsertReq) error {
	m, _ := s.FindOneByUserIdAndKey(userId, settings.Key)
	if m == nil {
		m = &model.Settings{
			UserId: userId,
			Key:    settings.Key,
			Value:  settings.Value,
		}
	} else {
		m.Value = settings.Value
	}
	return database.GetMysql().Save(m).Error
}

func (s *SettingsService) FindOneByUserIdAndKey(userId int64, key string) (*model.Settings, error) {
	m := &model.Settings{}
	err := database.GetMysql().Model(&model.Settings{}).Where("user_id = ? and `key` = ?", userId, key).First(m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, E.Message("查询的记录不存在")
		}
		return nil, err
	}
	return m, nil
}
