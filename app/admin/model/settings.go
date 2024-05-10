package model

type Settings struct {
	Id     int64  `json:"id" gorm:"type:bigint(20);primaryKey;autoIncrement;comment:主键ID"`
	UserId int64  `json:"userId" gorm:"type:bigint(20);not null;index:idx_uid_key,unique,priority:1;comment:用户ID"`
	Key    string `json:"key" gorm:"type:varchar(50);not null;index:idx_uid_key,unique,priority:2;comment:键名"`
	Value  string `json:"value" gorm:"type:varchar(2048);comment:值"`
}

func (*Settings) TableName() string {
	return "ops_settings"
}
