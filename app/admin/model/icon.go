package model

type Icon struct {
	Id    int64  `json:"id" gorm:"primaryKey;type:bigint(20);autoIncrement;comment:主键ID"`
	Value string `json:"value" gorm:"type:varchar(255);not null;comment:图标"`
}

func (*Icon) TableName() string {
	return "ops_icon"
}
