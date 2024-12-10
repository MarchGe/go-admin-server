package state

import (
	"context"
	"github.com/MarchGe/go-admin-server/app/admin/model/dvmodel/task"
	"github.com/MarchGe/go-admin-server/app/admin/service/dvservice/dto/req"
)

var _ TState = (*ActivatedState)(nil)
var _activatedState = &ActivatedState{
	tsc: GetTaskStateCommon(),
}

type ActivatedState struct {
	tsc *TaskStateCommon
}

func GetActivatedState() *ActivatedState {
	return _activatedState
}

func (s *ActivatedState) Start(ctx context.Context, t *task.ScriptTask) error {
	return UnsupportedOperationErr
}

func (s *ActivatedState) Run(ctx context.Context, t *task.ScriptTask) error {
	return s.tsc.Run(ctx, t)
}

func (s *ActivatedState) Stop(ctx context.Context, t *task.ScriptTask) error {
	return s.tsc.Stop(ctx, t)
}

func (s *ActivatedState) Update(ctx context.Context, info *req.ScriptTaskUpsertReq, t *task.ScriptTask) error {
	return UnsupportedOperationErr
}

func (s *ActivatedState) Delete(ctx context.Context, t *task.ScriptTask) error {
	return UnsupportedOperationErr
}
