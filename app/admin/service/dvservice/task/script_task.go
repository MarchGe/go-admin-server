package task

import (
	"context"
	"errors"
	"fmt"
	"github.com/MarchGe/go-admin-server/app/admin/model/dvmodel/task"
	"github.com/MarchGe/go-admin-server/app/admin/scheduler"
	"github.com/MarchGe/go-admin-server/app/admin/service/dto/res"
	"github.com/MarchGe/go-admin-server/app/admin/service/dvservice/dto/req"
	"github.com/MarchGe/go-admin-server/app/admin/service/dvservice/task/state"
	"github.com/MarchGe/go-admin-server/app/common/E"
	"github.com/MarchGe/go-admin-server/app/common/database"
	"gorm.io/gorm"
	"log/slog"
	"strings"
	"time"
)

var _scriptTaskService = &ScriptTaskService{}

type ScriptTaskService struct {
}

func GetScriptTaskService() *ScriptTaskService {
	return _scriptTaskService
}

func InitActivatedTasks() {
	tasks := make([]*task.ScriptTask, 0)
	err := database.GetMysql().Model(&task.ScriptTask{}).Preload("Scripts").Preload("HostGroup").Preload("HostGroup.HostList").
		Where("status = ? and execute_type = ?", task.StatusActive, task.ExecuteTypeAuto).
		Find(&tasks).Error
	if err != nil {
		slog.Error("Initiating activated script task to scheduler error", slog.Any("err", err))
	}
	for _, t := range tasks {
		scheduler.AddTask(t.Id, t.Cron, func() {
			// uId直接传递0，因为自动执行的脚本任务不需要用户id
			if err = GetScriptTaskService().Run(0, t); err != nil {
				slog.Error("Scheduler execute script task error", slog.Int64("taskId", t.Id), slog.Any("err", err))
			}
		})
		slog.Info("Add script task to scheduler success", slog.Int64("id", t.Id), slog.String("name", t.Name), slog.String("cron", t.Cron))
	}
}

func (s *ScriptTaskService) Create(info *req.ScriptTaskUpsertReq) error {
	t := s.toModel(info)
	t.Status = task.StatusNotRunning
	t.CreateTime = time.Now()
	t.UpdateTime = time.Now()
	err := database.GetMysql().Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(t).Error; err != nil {
			return err
		}
		if err := tx.Where("task_id = ?", t.Id).Delete(&task.ScriptTaskScript{}).Error; err != nil {
			return err
		}
		scripts := make([]*task.ScriptTaskScript, len(info.ScriptIds))
		for i := range info.ScriptIds {
			scripts[i] = &task.ScriptTaskScript{
				TaskId:   t.Id,
				ScriptId: info.ScriptIds[i],
			}
		}
		return tx.Save(scripts).Error
	})
	return err
}

func (s *ScriptTaskService) toModel(info *req.ScriptTaskUpsertReq) *task.ScriptTask {
	t := &task.ScriptTask{
		Name: info.Name,
		Kind: info.Kind,
		Cron: info.Cron,
	}
	if strings.TrimSpace(info.Cron) == "" {
		t.ExecuteType = task.ExecuteTypeManual
	} else {
		t.ExecuteType = task.ExecuteTypeAuto
	}
	if t.Kind == task.KindRemote {
		t.HostGroupId = info.HostGroupId
	}
	return t
}

func (s *ScriptTaskService) Update(id int64, info *req.ScriptTaskUpsertReq) error {
	t, err := s.FindOneById(id)
	if err != nil {
		return err
	}
	tState, err := s.getState(t.Status)
	if err != nil {
		return err
	}
	return tState.Update(context.TODO(), info, t)
}

func (s *ScriptTaskService) FindOneById(id int64, preloads ...string) (*task.ScriptTask, error) {
	m := &task.ScriptTask{}
	db := database.GetMysql()
	if len(preloads) > 0 {
		for _, preload := range preloads {
			db = db.Preload(preload)
		}
	}
	err := db.First(m, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, E.Message("任务不存在")
		} else {
			return nil, err
		}
	}
	return m, nil
}

func (s *ScriptTaskService) Delete(id int64) error {
	t, err := s.FindOneById(id)
	if err != nil {
		return err
	}
	tState, err := s.getState(t.Status)
	if err != nil {
		return err
	}
	return tState.Delete(context.TODO(), t)
}

func (s *ScriptTaskService) PageList(keyword string, page, pageSize int) (*res.PageableData[*task.ScriptTask], error) {
	tasks := make([]*task.ScriptTask, 0)
	pageableData := &res.PageableData[*task.ScriptTask]{}
	db := database.GetMysql().Model(&task.ScriptTask{}).Preload("Scripts", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name", "version")
	})
	if keyword != "" {
		db.Where("name like ?", "%"+keyword+"%")
	}
	var count int64
	err := db.Count(&count).Order("create_time desc").Offset(pageSize * (page - 1)).Limit(pageSize).Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	pageableData.List = tasks
	pageableData.Total = count
	return pageableData, nil
}

func (s *ScriptTaskService) Start(uId, id int64) error {
	t, err := s.FindOneById(id, "Scripts", "HostGroup", "HostGroup.HostList")
	if err != nil {
		return err
	}
	tState, err := s.getState(t.Status)
	if err != nil {
		return err
	}
	ctx := context.WithValue(context.Background(), "uId", uId)
	return tState.Start(ctx, t)
}

func (s *ScriptTaskService) Run(uId int64, t *task.ScriptTask) error {
	tState, err := s.getState(t.Status)
	if err != nil {
		return err
	}
	ctx := context.WithValue(context.Background(), "uId", uId)
	return tState.Run(ctx, t)
}

func (s *ScriptTaskService) Stop(id int64) error {
	t, err := s.FindOneById(id)
	if err != nil {
		return err
	}
	tState, err := s.getState(t.Status)
	if err != nil {
		return err
	}
	return tState.Stop(context.TODO(), t)
}

func (s *ScriptTaskService) getState(status task.Status) (state.TState, error) {
	switch status {
	case task.StatusNotRunning:
		return state.GetNotRunningState(), nil
	case task.StatusRunning:
		return state.GetRunningState(), nil
	case task.StatusComplete:
		return state.GetCompletedState(), nil
	case task.StatusStopped:
		return state.GetStoppedState(), nil
	case task.StatusActive:
		return state.GetActivatedState(), nil
	default:
		return nil, fmt.Errorf("unknown script task status: %v", status)
	}
}
