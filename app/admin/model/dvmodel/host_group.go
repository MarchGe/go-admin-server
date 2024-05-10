package dvmodel

type HostGroup struct {
	Id      int64 `json:"id" gorm:"primaryKey;type:bigint(20);autoIncrement;comment:主键ID"`
	HostId  int64 `json:"hostId" gorm:"type:bigint(20);not null;index:idx_h_g_id,unique,priority:2;index:idx_h_id;comment:主机ID"`
	GroupId int64 `json:"groupId" gorm:"type:bigint(20);not null;index:idx_h_g_id,unique,priority:1;comment:分组ID"`
}

func (*HostGroup) TableName() string {
	return "dv_host_group"
}
