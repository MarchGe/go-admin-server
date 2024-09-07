package req

type ExplorerCreateDirReq struct {
	Dir  string `json:"dir" binding:"required,max=1024" label:"当前目录"`
	Name string `json:"name" binding:"required,max=50,regex=^[^<>:\"?*\x00-\x1f\x60\\x7c]+$" label:"目录名称"`
}
