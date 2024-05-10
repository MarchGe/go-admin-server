package model

type RoleMenu struct {
	Id     int64 `json:"id" gorm:"primaryKey;type:bigint(20);autoIncrement;comment:主键ID"`
	RoleId int64 `json:"roleId" gorm:"type:bigint(20);not null;index:idx_r_m_id,unique,priority:1;comment:角色ID"`
	MenuId int64 `json:"menuId" gorm:"type:bigint(20);not null;index:idx_r_m_id,unique,priority:2;index:idx_m_id;comment:菜单ID"`
}

func (*RoleMenu) TableName() string {
	return "ops_role_menu"
}
