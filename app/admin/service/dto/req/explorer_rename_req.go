package req

type ExplorerRenameReq struct {
	Dir     string `json:"dir" binding:"required,max=1024" label:"当前目录"`
	OldName string `json:"oldName" binding:"required,max=1024" label:"旧名称"`
	NewName string `json:"newName" binding:"required,max=50,regex=^[^<>:\"?*\x00-\x1f\x60\\x7c]+$" label:"新名称"`
}
