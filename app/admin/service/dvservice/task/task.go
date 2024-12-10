package task

import (
	"context"
	"errors"
	"github.com/MarchGe/go-admin-server/app/admin/model/dvmodel/task"
	"github.com/MarchGe/go-admin-server/app/admin/service/dto/res"
	"github.com/MarchGe/go-admin-server/app/admin/service/dvservice/dto/req"
	dvRes "github.com/MarchGe/go-admin-server/app/admin/service/dvservice/dto/res"
	"github.com/MarchGe/go-admin-server/app/common/E"
	"github.com/MarchGe/go-admin-server/app/common/database"
	"gorm.io/gorm"
	"time"
)

var _taskService = &TaskService{}

type TaskService struct {
}

func GetTaskService() *TaskService {
	return _taskService
}

func (s *TaskService) CreateTask(info *req.TaskUpsertReq) error {
	t := s.toModel(info)
	t.Status = task.StatusNotRunning
	t.CreateTime = time.Now()
	t.UpdateTime = time.Now()
	err := database.GetMysql().Transaction(func(tx *gorm.DB) error {
		taskType := t.Type
		id, err := Select(taskType).Create(tx, info.Concrete)
		if err != nil {
			return err
		}
		t.AssociationID = id
		t.AssociationType = GetPolymorphicValue(taskType)
		return tx.Save(t).Error
	})
	return err
}

func (s *TaskService) toModel(info *req.TaskUpsertReq) *task.Task {
	return &task.Task{
		Name:        info.Name,
		Type:        info.Type,
		Cron:        info.Cron,
		ExecuteType: info.ExecuteType,
	}
}

func (s *TaskService) UpdateTask(id int64, info *req.TaskUpsertReq) error {
	t, err := s.FindOneById(id)
	if err != nil {
		return err
	}
	if t.Status == task.StatusRunning {
		return E.Message("已启动的任务，不允许修改")
	}
	s.copyProperties(info, t)
	t.UpdateTime = time.Now()
	err = database.GetMysql().Transaction(func(tx *gorm.DB) error {
		taskType := t.Type
		if err = Select(taskType).Update(tx, info.Concrete, t.AssociationID); err != nil {
			return err
		}
		return tx.Save(t).Error
	})
	return err
}

func (s *TaskService) FindOneById(id int64) (*task.Task, error) {
	m := &task.Task{}
	err := database.GetMysql().First(m, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, E.Message("任务不存在")
		} else {
			return nil, err
		}
	}
	return m, nil
}

func (s *TaskService) copyProperties(info *req.TaskUpsertReq, task *task.Task) {
	task.Name = info.Name
	task.Type = info.Type
	task.Cron = info.Cron
	task.ExecuteType = info.ExecuteType
}

func (s *TaskService) DeleteTask(id int64) error {
	t, err := s.FindOneById(id)
	if err != nil {
		return err
	}
	err = database.GetMysql().Transaction(func(tx *gorm.DB) error {
		if err = Select(t.Type).Delete(tx, t); err != nil {
			return err
		}
		return tx.Delete(&task.Task{}, id).Error
	})
	return err
}

func (s *TaskService) PageList(keyword string, page, pageSize int) (*res.PageableData[*dvRes.TaskRes], error) {
	tasks := make([]*task.Task, 0)
	pageableData := &res.PageableData[*dvRes.TaskRes]{}
	db := database.GetMysql().Model(&task.Task{})
	if keyword != "" {
		db.Where("name like ?", "%"+keyword+"%")
	}
	var count int64
	err := db.Count(&count).Order("create_time desc").Offset(pageSize * (page - 1)).Limit(pageSize).Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	resTasks := s.attachConcreteTask(tasks)
	pageableData.List = resTasks
	pageableData.Total = count
	return pageableData, nil
}

func (s *TaskService) attachConcreteTask(tasks []*task.Task) []*dvRes.TaskRes {
	taskRes := make([]*dvRes.TaskRes, len(tasks))
	for i := range tasks {
		t := &dvRes.TaskRes{
			Task: tasks[i],
		}
		t.Concrete, _ = Select(tasks[i].Type).FindOneById(tasks[i].AssociationID)
		taskRes[i] = t
	}
	return taskRes
}

func (s *TaskService) StartTask(uId, id int64) error {
	t, err := s.FindOneById(id)
	if err != nil {
		return err
	}
	if t.Status == task.StatusRunning {
		return E.Message("当前任务已经启动了")
	}
	parentCtx := context.WithValue(context.Background(), "uId", uId)
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()
	return Select(t.Type).Start(ctx, t)
}

func (s *TaskService) StopTask(id int64) error {
	t, err := s.FindOneById(id)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	return Select(t.Type).Stop(ctx, t)
}
