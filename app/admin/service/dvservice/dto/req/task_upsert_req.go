package req

type TaskUpsertReq struct {
	Name        string         `json:"name" binding:"required,max=50" label:"名称"`
	Type        int8           `json:"type" binding:"required,oneof=1" label:"任务类型"`
	Cron        string         `json:"cron" binding:"omitempty,max=50" label:"cron表达式"`
	ExecuteType int8           `json:"executeType" binding:"omitempty,oneof=0 1" label:"执行方式"`
	Concrete    map[string]any `json:"concrete" binding:"required" label:"具体任务"`
}
