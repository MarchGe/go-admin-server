package dvmodel

import "github.com/MarchGe/go-admin-server/app/admin/model"

type Group struct {
	model.Base
	Name     string  `json:"name" gorm:"type:varchar(50);not null;index:idx_name;comment:分组名称"`
	SortNum  int32   `json:"sortNum" gorm:"type:int;default:0;comment:顺序"`
	HostList []*Host `json:"hostList" gorm:"-:migration;many2many:dv_host_group"`
}

func (*Group) TableName() string {
	return "dv_group"
}
