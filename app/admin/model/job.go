package model

type Job struct {
	Base
	Name        string `json:"name" gorm:"type:varchar(50);not null;index:idx_name;comment:岗位名称"`
	SortNum     int    `json:"sortNum" gorm:"type:int;default:0;comment:排序"`
	Description string `json:"description" gorm:"type:varchar(255);comment:岗位描述"`
}

func (*Job) TableName() string {
	return "ops_job"
}
