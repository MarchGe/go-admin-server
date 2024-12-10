package state

import (
	"context"
	"github.com/MarchGe/go-admin-server/app/admin/model/dvmodel/task"
	"github.com/MarchGe/go-admin-server/app/admin/service/dvservice/dto/req"
)

var _ TState = (*NotRunningState)(nil)
var _notRunningState = &NotRunningState{
	tsc: GetTaskStateCommon(),
}

type NotRunningState struct {
	tsc *TaskStateCommon
}

func GetNotRunningState() *NotRunningState {
	return _notRunningState
}

func (s *NotRunningState) Start(ctx context.Context, t *task.ScriptTask) error {
	return s.tsc.Start(ctx, t)
}

func (s *NotRunningState) Run(ctx context.Context, t *task.ScriptTask) error {
	return s.tsc.Run(ctx, t)
}

func (s *NotRunningState) Stop(ctx context.Context, t *task.ScriptTask) error {
	return UnsupportedOperationErr
}

func (s *NotRunningState) Update(ctx context.Context, info *req.ScriptTaskUpsertReq, t *task.ScriptTask) error {
	return s.tsc.Update(ctx, info, t)
}

func (s *NotRunningState) Delete(ctx context.Context, t *task.ScriptTask) error {
	return s.tsc.Delete(ctx, t)
}
