package state

import (
	"context"
	"github.com/MarchGe/go-admin-server/app/admin/model/dvmodel/task"
	"github.com/MarchGe/go-admin-server/app/admin/service/dvservice/dto/req"
)

var _ TState = (*RunningState)(nil)
var _runningState = &RunningState{
	tsc: GetTaskStateCommon(),
}

type RunningState struct {
	tsc *TaskStateCommon
}

func GetRunningState() *RunningState {
	return _runningState
}

func (s *RunningState) Start(ctx context.Context, t *task.ScriptTask) error {
	return UnsupportedOperationErr
}

func (s *RunningState) Run(ctx context.Context, t *task.ScriptTask) error {
	return UnsupportedOperationErr
}

func (s *RunningState) Stop(ctx context.Context, t *task.ScriptTask) error {
	return s.tsc.Stop(ctx, t)
}

func (s *RunningState) Update(ctx context.Context, info *req.ScriptTaskUpsertReq, t *task.ScriptTask) error {
	return UnsupportedOperationErr
}

func (s *RunningState) Delete(ctx context.Context, t *task.ScriptTask) error {
	return UnsupportedOperationErr
}
