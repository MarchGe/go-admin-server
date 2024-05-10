package req

type JobUpsertReq struct {
	Name        string `json:"name" binding:"required,max=50" label:"岗位名称"`
	SortNum     int    `json:"sortNum" binding:"omitempty" label:"排序"`
	Description string `json:"description" binding:"omitempty,max=255" label:"岗位名称"`
}
