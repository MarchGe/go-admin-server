package req

type ScriptUpsertReq struct {
	Name        string `json:"name" binding:"required,max=50" label:"名称"`
	Version     string `json:"version" binding:"required,max=50" label:"版本"`
	Content     string `json:"content" binding:"required,max=10000" label:"脚本内容"`
	Description string `json:"description" binding:"omitempty,max=2000" label:"使用说明"`
}
