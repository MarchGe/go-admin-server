package task

type DeployTask struct {
	Id          int64  `json:"id" gorm:"type:bigint(20);primaryKey;autoIncrement;comment:主键ID"`
	UploadPath  string `json:"uploadPath" gorm:"type:varchar(255);not null;comment:部署包的上传路劲"`
	AppId       int64  `json:"appId" gorm:"type:bigint(20);not null;comment:关联的应用ID"`
	ScriptId    int64  `json:"scriptId" gorm:"type:bigint(20);not null;comment:关联的部署脚本ID"`
	HostGroupId int64  `json:"hostGroupId" gorm:"type:bigint(20);not null;comment:关联的服务器分组ID"`
	Task        *Task  `json:"task" gorm:"polymorphic:Association;polymorphicValue:deploy"` // 多态
}

func (*DeployTask) TableName() string {
	return "dv_task_deploy"
}
