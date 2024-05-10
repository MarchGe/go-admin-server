package dvmodel

import "github.com/MarchGe/go-admin-server/app/admin/model"

type Host struct {
	model.Base
	Name     string `json:"name" gorm:"type:varchar(50);not null;index:idx_name;comment:名称"`
	Ip       string `json:"ip" gorm:"type:varchar(50);not null;index:idx_ip,unique;comment:IP"`
	Port     int16  `json:"port" gorm:"type:int;not null;comment:ssh端口"`
	User     string `json:"user" gorm:"type:varchar(50);not null;comment:ssh用户名"`
	Password string `json:"password" gorm:"type:varchar(255);not null;comment:ssh密码"`
	SortNum  int32  `json:"sortNum" gorm:"type:int;default:0;comment:顺序"`
}

func (*Host) TableName() string {
	return "dv_host"
}
