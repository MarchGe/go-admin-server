package model

import "time"

type Base struct {
	Id         int64     `json:"id" gorm:"type:bigint(20);primaryKey;autoIncrement;comment:主键ID"`
	CreateTime time.Time `json:"createTime" gorm:"type:datetime;not null;comment:创建时间"`
	UpdateTime time.Time `json:"updateTime" gorm:"type:datetime;default:CURRENT_TIMESTAMP;autoUpdateTime;comment:更新时间"`
}
