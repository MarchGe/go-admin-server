package req

import (
	"github.com/MarchGe/go-admin-server/app/admin/model/dvmodel/task"
	"github.com/MarchGe/go-admin-server/app/admin/scheduler"
	"github.com/MarchGe/go-admin-server/app/common/E"
)

type ScriptTaskUpsertReq struct {
	Name        string    `json:"name" binding:"required,max=50" label:"名称"`
	Cron        string    `json:"cron" binding:"omitempty,max=50" label:"cron表达式"`
	Kind        task.Kind `json:"kind" binding:"omitempty,oneof=0 1" label:"任务类型"`
	ScriptIds   []int64   `json:"scriptIds" binding:"required,min=1" label:"脚本列表"`
	HostGroupId int64     `json:"hostGroupId" binding:"omitempty,min=1" label:"服务器分组ID"`
}

func (s *ScriptTaskUpsertReq) Verify() error {
	for _, item := range s.ScriptIds {
		if item <= 0 {
			return E.Message("脚本项不能为空")
		}
	}
	if s.Kind == task.KindRemote {
		if s.HostGroupId <= 0 {
			return E.Message("服务器分组不能为空")
		}
	}
	if s.Cron != "" {
		if err := scheduler.VerifyCron(s.Cron); err != nil {
			return err
		}
	}
	return nil
}
