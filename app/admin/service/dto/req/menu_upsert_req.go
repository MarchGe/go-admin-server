package req

type MenuUpsertReq struct {
	Symbol      string `json:"symbol" binding:"omitempty,max=50,regex=^\\w+:\\w+$" label:"权限标识"`
	Name        string `json:"name" binding:"required,max=20" label:"菜单名称"`
	Icon        string `json:"icon" binding:"omitempty,max=50" label:"图标"`
	SortNum     int32  `json:"sortNum" binding:"omitempty" label:"排序"`
	Url         string `json:"url" binding:"omitempty,max=100" label:"菜单URL"`
	Display     int8   `json:"display" binding:"omitempty,oneof=0 1" label:"是否显示"`
	External    int8   `json:"external" binding:"omitempty,oneof=0 1" label:"是否外链"`
	ExternalWay int8   `json:"externalWay" binding:"omitempty,oneof=0 1" label:"打开方式"`
	ParentId    int64  `json:"parentId" binding:"omitempty,min=0" label:"父菜单"`
}
