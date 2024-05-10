package dvmodel

import "github.com/MarchGe/go-admin-server/app/admin/model"

type Script struct {
	model.Base
	Name        string `json:"name" gorm:"type:varchar(50);not null;index:idx_n_v,unique,priority:1;comment:名称"`
	Version     string `json:"version" gorm:"type:varchar(50);not null;index:idx_n_v,unique,priority:2;comment:版本"`
	Content     string `json:"content" gorm:"type:varchar(10000);not null;comment:脚本内容"`
	Description string `json:"description" gorm:"type:varchar(2000);not null;comment:使用说明"`
}

func (*Script) TableName() string {
	return "dv_script"
}
