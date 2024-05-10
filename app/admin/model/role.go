package model

type Role struct {
	Base
	Name     string  `json:"name" gorm:"type:varchar(50);not null;index:idx_name;comment:角色名称"`
	SortNum  int     `json:"sortNum" gorm:"type:int;default:0;comment:排序"`
	MenuList []*Menu `json:"menuList" gorm:"-:migration;many2many:ops_role_menu"`
}

func (*Role) TableName() string {
	return "ops_role"
}
