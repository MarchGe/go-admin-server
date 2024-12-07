package state

import (
	"context"
	"github.com/MarchGe/go-admin-server/app/admin/model/dvmodel/task"
	"github.com/MarchGe/go-admin-server/app/admin/service/dvservice/dto/req"
)

var _ TState = (*StoppedState)(nil)
var _stoppedState = &StoppedState{
	tsc: GetTaskStateCommon(),
}

type StoppedState struct {
	tsc *TaskStateCommon
}

func GetStoppedState() *StoppedState {
	return _stoppedState
}

func (s *StoppedState) Start(ctx context.Context, t *task.ScriptTask) error {
	return s.tsc.Start(ctx, t)
}

func (s *StoppedState) Run(ctx context.Context, t *task.ScriptTask) error {
	return s.tsc.Run(ctx, t)
}

func (s *StoppedState) Stop(ctx context.Context, t *task.ScriptTask) error {
	return UnsupportedOperationErr
}

func (s *StoppedState) Update(ctx context.Context, info *req.ScriptTaskUpsertReq, t *task.ScriptTask) error {
	return s.tsc.Update(ctx, info, t)
}

func (s *StoppedState) Delete(ctx context.Context, t *task.ScriptTask) error {
	return s.tsc.Delete(ctx, t)
}
