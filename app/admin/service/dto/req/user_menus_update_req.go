package req

type UserMenusUpdateReq struct {
	Ids []int64 `json:"ids" binding:"required,max=1000" label:"菜单列表"`
}
