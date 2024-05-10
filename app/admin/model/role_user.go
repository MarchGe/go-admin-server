package model

type RoleUser struct {
	Id     int64 `json:"id" gorm:"primaryKey;type:bigint(20);autoIncrement;comment:主键ID"`
	RoleId int64 `json:"roleId" gorm:"type:bigint(20);not null;index:idx_r_u_id,unique,priority:1;comment:角色ID"`
	UserId int64 `json:"userId" gorm:"type:bigint(20);not null;index:idx_r_u_id,unique,priority:2;index:idx_uid;comment:用户ID"`
}

func (*RoleUser) TableName() string {
	return "ops_role_user"
}
