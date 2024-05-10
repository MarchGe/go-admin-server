package task

import (
	"context"
	"errors"
	"fmt"
	"github.com/MarchGe/go-admin-server/app/admin/model/dvmodel"
	"github.com/MarchGe/go-admin-server/app/common/E"
	"gorm.io/gorm"
)

type ConcreteTask interface {
	FindOneById(id int64) (*dvmodel.DeployTask, error)
	Create(tx *gorm.DB, data map[string]any) (id int64, err error)
	Update(tx *gorm.DB, data map[string]any, id int64) error
	Delete(tx *gorm.DB, t *dvmodel.Task) error
	Start(ctx context.Context, t *dvmodel.Task) error
	Run(ctx context.Context, t *dvmodel.Task) error
	Stop(ctx context.Context, t *dvmodel.Task) error
}

func Select(taskType dvmodel.TaskType) ConcreteTask {
	switch taskType {
	case dvmodel.TaskTypeDeploy:
		return _deployTaskService
	default:
		E.PanicErr(errors.New(fmt.Sprintf("Unknown task type: %d", taskType)))
		return nil
	}
}

func GetPolymorphicValue(taskType dvmodel.TaskType) string {
	switch taskType {
	case dvmodel.TaskTypeDeploy:
		return "deploy"
	default:
		E.PanicErr(errors.New(fmt.Sprintf("Unknown task type: %d", taskType)))
		return ""
	}
}
