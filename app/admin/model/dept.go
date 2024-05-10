package model

type Dept struct {
	Base
	Name     string `json:"name" gorm:"type:varchar(50);not null;index:idx_name;comment:部门名称"`
	SortNum  int32  `json:"sortNum" gorm:"type:int;default:0;comment:部门顺序"`
	ParentId int64  `json:"parentId" gorm:"type:bigint(20);default:0;index:idx_pid;comment:父级部门ID"`
}

func (*Dept) TableName() string {
	return "ops_dept"
}
