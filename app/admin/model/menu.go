package model

type Menu struct {
	Base
	Symbol      string `json:"symbol" gorm:"type:varchar(100);index:idx_symbol;comment:权限标识"`
	Name        string `json:"name" gorm:"type:varchar(50);not null;index:idx_name;comment:菜单名称"`
	Icon        string `json:"icon" gorm:"type:varchar(255);comment:菜单图标"`
	SortNum     int32  `json:"sortNum" gorm:"type:int;default:0;comment:菜单顺序"`
	Url         string `json:"url" gorm:"type:varchar(255);comment:路由路径"`
	Display     int8   `json:"display" gorm:"type:tinyint(4);default:0;comment:是否显示，0-否，1-是"`
	External    int8   `json:"external" gorm:"type:tinyint(4);default:0;comment:是否外链，0-否，1-是"`
	ExternalWay int8   `json:"externalWay" gorm:"type:tinyint(4);default:0;comment:外链打开方式（仅外链有效），0-外联，1-内嵌"`
	ParentId    int64  `json:"parentId" gorm:"type:bigint(20);default:0;index:idx_pid;comment:父菜单ID"`
}

func (*Menu) TableName() string {
	return "ops_menu"
}
