package task

import (
	"github.com/MarchGe/go-admin-server/app/admin/model"
	"github.com/MarchGe/go-admin-server/app/admin/model/dvmodel"
)

type ScriptTask struct {
	model.Base
	Name        string            `json:"name" gorm:"type:varchar(50);not null;index:idx_name;comment:名称"`
	Status      Status            `json:"status" gorm:"type:tinyint(4);default:0;comment:任务状态"`
	ExecuteType ExecuteType       `json:"executeType" gorm:"type:tinyint(4);default:0;comment:执行方式"`
	Scripts     []*dvmodel.Script `json:"scripts" gorm:"-:migration;many2many:dv_script_task_script;joinForeignKey:TaskId"`
	Kind        Kind              `json:"kind" gorm:"type:tinyint(4);default:0;comment:任务类型"`
	HostGroupId int64             `json:"hostGroupId" gorm:"type:bigint(20);comment:关联的服务器分组ID"`
	HostGroup   *dvmodel.Group    `json:"hostGroup" gorm:"-:migration;foreignKey:HostGroupId"`
	Cron        string            `json:"cron" gorm:"type:varchar(50);comment:cron表达式"`
}

func (*ScriptTask) TableName() string {
	return "dv_script_task"
}
