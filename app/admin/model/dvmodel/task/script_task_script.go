package task

type ScriptTaskScript struct {
	Id       int64 `json:"id" gorm:"type:bigint(20);primaryKey;autoIncrement;comment:主键ID"`
	TaskId   int64 `json:"taskId" gorm:"type:bigint(20);not null;index:idx_t_s_id,unique,priority:1;comment:脚本任务ID"`
	ScriptId int64 `json:"scriptId" gorm:"type:bigint(20);not null;index:idx_t_s_id,unique,priority:2;index:idx_s_id;comment:脚本ID"`
}

func (*ScriptTaskScript) TableName() string {
	return "dv_script_task_script"
}
