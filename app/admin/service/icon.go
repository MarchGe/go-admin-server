package service

import (
	"github.com/MarchGe/go-admin-server/app/admin/model"
	"github.com/MarchGe/go-admin-server/app/common/database"
)

var _iconService = &IconService{}

type IconService struct {
}

func GetIconService() *IconService {
	return _iconService
}

func (s *IconService) FindAll() ([]*model.Icon, error) {
	icons := make([]*model.Icon, 0)
	err := database.GetMysql().Model(&model.Icon{}).Find(&icons).Error
	if err != nil {
		return nil, err
	}
	return icons, nil
}
