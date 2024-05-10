package model

type UserMenu struct {
	Id     int64 `json:"id" gorm:"primaryKey;type:bigint(20);autoIncrement;comment:主键ID"`
	UserId int64 `json:"userId" gorm:"type:bigint(20);not null;index:idx_u_m_id,unique,priority:1;comment:用户ID"`
	MenuId int64 `json:"menuId" gorm:"type:bigint(20);not null;index:idx_u_m_id,unique,priority:2;index:idx_m_id;comment:菜单ID"`
}

func (*UserMenu) TableName() string {
	return "ops_user_menu"
}
