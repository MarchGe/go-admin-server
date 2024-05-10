package model

type UserPassword struct {
	Id       int64  `json:"id" gorm:"primaryKey;type:bigint(20);autoIncrement;comment:主键ID"`
	UserId   int64  `json:"userId" gorm:"type:bigint(20);not null;index:idx_uid,unique;comment:用户ID"`
	Password string `json:"password" gorm:"type:varchar(255);not null;comment:用户密码"`
}

func (*UserPassword) TableName() string {
	return "ops_user_password"
}
