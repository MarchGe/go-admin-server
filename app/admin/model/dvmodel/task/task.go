package task

import "github.com/MarchGe/go-admin-server/app/admin/model"

type Type int8

const (
	TypeDeploy Type = 1 // 任务类型 - 部署任务
)

type Task struct {
	model.Base
	Name            string      `json:"name" gorm:"type:varchar(50);not null;index:idx_name;comment:名称"`
	Type            Type        `json:"type" gorm:"type:tinyint(4);not null;comment:任务类型"`
	Status          Status      `json:"status" gorm:"type:tinyint(4);default:0;comment:任务状态"`
	AssociationID   int64       `json:"associationId" gorm:"type:bigint(20);not null;comment:关联的具体任务的ID"`
	AssociationType string      `json:"associationType" gorm:"type:varchar(50);not null;comment:关联的具体任务表类型"`
	Cron            string      `json:"cron" gorm:"type:varchar(50);comment:cron表达式"`
	ExecuteType     ExecuteType `json:"executeType" gorm:"type:tinyint(4);default:0;comment:执行方式"`
}

func (*Task) TableName() string {
	return "dv_task"
}
