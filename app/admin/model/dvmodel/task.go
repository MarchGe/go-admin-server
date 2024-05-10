package dvmodel

import "github.com/MarchGe/go-admin-server/app/admin/model"

type TaskType int8
type TaskStatus int8
type ExecuteType int8

const (
	TaskStatusNotRunning TaskStatus = 0 // 任务状态 - 未运行
	TaskStatusRunning    TaskStatus = 1 // 任务状态 - 运行中
	TaskStatusComplete   TaskStatus = 2 // 任务状态 - 已完成
	TaskStatusStopped    TaskStatus = 3 // 任务状态 - 已停止
)

const (
	TaskTypeDeploy TaskType = 1 // 任务类型 - 部署任务
)

const (
	ExecuteTypeManual ExecuteType = 0 // 执行方式 - 手动执行
	ExecuteTypeAuto   ExecuteType = 1 // 执行方式 - 自动执行，根据cron表达式
)

type Task struct {
	model.Base
	Name            string `json:"name" gorm:"type:varchar(50);not null;index:idx_name;comment:名称"`
	Type            int8   `json:"type" gorm:"type:tinyint(4);not null;comment:任务类型"`
	Status          int8   `json:"status" gorm:"type:tinyint(4);default:0;comment:任务状态"`
	AssociationID   int64  `json:"associationId" gorm:"type:bigint(20);not null;comment:关联的具体任务的ID"`
	AssociationType string `json:"associationType" gorm:"type:varchar(50);not null;comment:关联的具体任务表类型"`
	Cron            string `json:"cron" gorm:"type:varchar(50);comment:cron表达式"`
	ExecuteType     int8   `json:"executeType" gorm:"type:tinyint(4);default:0;comment:执行方式"`
}

func (*Task) TableName() string {
	return "dv_task"
}
