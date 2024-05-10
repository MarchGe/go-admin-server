package dvmodel

import "github.com/MarchGe/go-admin-server/app/admin/model"

type App struct {
	model.Base
	Name     string `json:"name" gorm:"type:varchar(50);not null;index:idx_n_v,unique,priority:1;comment:名称"`
	Version  string `json:"version" gorm:"type:varchar(50);not null;index:idx_n_v,unique,priority:2;comment:版本"`
	Port     int16  `json:"port" gorm:"type:int;not null;comment:应用端口"`
	Key      string `json:"key" gorm:"type:varchar(255);not null;comment:部署包的路劲"`
	FileName string `json:"fileName" gorm:"type:varchar(100);not null;comment:部署包的文件名"`
}

func (*App) TableName() string {
	return "dv_app"
}
