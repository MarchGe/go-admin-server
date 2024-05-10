package req

type DeptUpsertReq struct {
	Name     string `json:"name" binding:"required,max=50" label:"部门名称"`
	SortNum  int32  `json:"sortNum" binding:"omitempty" label:"排序"`
	ParentId int64  `json:"parentId" binding:"omitempty,min=0" label:"父级部门"`
}
