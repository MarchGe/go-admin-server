package model

type UserJob struct {
	Id     int64 `json:"id" gorm:"primaryKey;type:bigint(20);autoIncrement;comment:主键ID"`
	UserId int64 `json:"userId" gorm:"type:bigint(20);not null;index:idx_u_j_id,unique,priority:1;comment:用户ID"`
	JobId  int64 `json:"jobId" gorm:"type:bigint(20);not null;index:idx_u_j_id,unique,priority:2;index:idx_j_id;comment:岗位ID"`
}

func (*UserJob) TableName() string {
	return "ops_user_job"
}
