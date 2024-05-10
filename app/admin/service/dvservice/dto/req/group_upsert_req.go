package req

type GroupUpsertReq struct {
	Name    string  `json:"name" binding:"required,max=50" label:"分组名称"`
	SortNum int32   `json:"sortNum" binding:"omitempty" label:"排序"`
	HostIds []int64 `json:"hostIds" binding:"required,min=1,max=10000" label:"服务器列表"`
}
