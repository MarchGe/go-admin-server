package state

import (
	"context"
	"github.com/MarchGe/go-admin-server/app/admin/model/dvmodel/task"
	"github.com/MarchGe/go-admin-server/app/admin/service/dvservice/dto/req"
)

type TState interface {
	Start(ctx context.Context, t *task.ScriptTask) error
	Run(ctx context.Context, t *task.ScriptTask) error
	Stop(ctx context.Context, t *task.ScriptTask) error
	Update(ctx context.Context, info *req.ScriptTaskUpsertReq, t *task.ScriptTask) error
	Delete(ctx context.Context, t *task.ScriptTask) error
}
