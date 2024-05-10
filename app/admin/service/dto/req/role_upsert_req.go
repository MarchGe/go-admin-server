package req

type RoleUpsertReq struct {
	Name    string `json:"name" binding:"required,max=20" label:"角色名称"`
	SortNum int    `json:"sortNum" binding:"omitempty" label:"排序"`
}
