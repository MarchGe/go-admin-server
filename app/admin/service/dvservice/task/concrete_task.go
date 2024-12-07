package task

import (
	"context"
	"errors"
	"fmt"
	"github.com/MarchGe/go-admin-server/app/admin/model/dvmodel/task"
	"github.com/MarchGe/go-admin-server/app/common/E"
	"gorm.io/gorm"
)

type ConcreteTask interface {
	FindOneById(id int64) (*task.DeployTask, error)
	Create(tx *gorm.DB, data map[string]any) (id int64, err error)
	Update(tx *gorm.DB, data map[string]any, id int64) error
	Delete(tx *gorm.DB, t *task.Task) error
	Start(ctx context.Context, t *task.Task) error
	Run(ctx context.Context, t *task.Task) error
	Stop(ctx context.Context, t *task.Task) error
}

func Select(taskType task.Type) ConcreteTask {
	switch taskType {
	case task.TypeDeploy:
		return _deployTaskService
	default:
		E.PanicErr(errors.New(fmt.Sprintf("Unknown task type: %d", taskType)))
		return nil
	}
}

func GetPolymorphicValue(taskType task.Type) string {
	switch taskType {
	case task.TypeDeploy:
		return "deploy"
	default:
		E.PanicErr(errors.New(fmt.Sprintf("Unknown task type: %d", taskType)))
		return ""
	}
}
