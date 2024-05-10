package req

type ApiUpsertReq struct {
	Method      string `json:"method" binding:"required,max=10" label:"请求方式"`
	Name        string `json:"name" binding:"omitempty,max=20" label:"API名称"`
	Path        string `json:"path" binding:"required,max=100" label:"路由路径"`
	Description string `json:"description" binding:"omitempty,max=255" label:"描述信息"`
}
