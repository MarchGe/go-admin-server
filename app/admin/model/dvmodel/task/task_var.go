package task

type Status int8
type ExecuteType int8
type Kind int8

const (
	StatusNotRunning Status = 0 // 任务运行状态 - 未运行
	StatusRunning    Status = 1 // 任务运行状态 - 运行中
	StatusComplete   Status = 2 // 任务运行状态 - 已完成
	StatusStopped    Status = 3 // 任务运行状态 - 失败
	StatusActive     Status = 4 // 任务运行状态 - 已激活 （定时任务特有的状态）
)

const (
	ExecuteTypeManual ExecuteType = 0 // 任务执行方式 - 手动执行
	ExecuteTypeAuto   ExecuteType = 1 // 任务执行方式 - 自动执行，根据cron表达式
)

const (
	KindLocal  Kind = 0 // 任务类型 - 本地任务
	KindRemote Kind = 1 // 任务类型 - 远程任务
)
