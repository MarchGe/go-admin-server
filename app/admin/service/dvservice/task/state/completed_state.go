package state

import (
	"context"
	"github.com/MarchGe/go-admin-server/app/admin/model/dvmodel/task"
	"github.com/MarchGe/go-admin-server/app/admin/service/dvservice/dto/req"
)

var _ TState = (*CompletedState)(nil)
var _completeState = &CompletedState{
	tsc: GetTaskStateCommon(),
}

type CompletedState struct {
	tsc *TaskStateCommon
}

func GetCompletedState() *CompletedState {
	return _completeState
}

func (s *CompletedState) Start(ctx context.Context, t *task.ScriptTask) error {
	return s.tsc.Start(ctx, t)
}

func (s *CompletedState) Run(ctx context.Context, t *task.ScriptTask) error {
	return s.tsc.Run(ctx, t)
}

func (s *CompletedState) Stop(ctx context.Context, t *task.ScriptTask) error {
	return UnsupportedOperationErr
}

func (s *CompletedState) Update(ctx context.Context, info *req.ScriptTaskUpsertReq, t *task.ScriptTask) error {
	return s.tsc.Update(ctx, info, t)
}

func (s *CompletedState) Delete(ctx context.Context, t *task.ScriptTask) error {
	return s.tsc.Delete(ctx, t)
}
