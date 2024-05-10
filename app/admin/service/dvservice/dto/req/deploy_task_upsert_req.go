package req

type DeployTaskUpsertReq struct {
	UploadPath  string `json:"uploadPath" binding:"required,max=255" label:"上传路径"`
	AppId       int64  `json:"appId" binding:"required,min=1" label:"应用"`
	ScriptId    int64  `json:"scriptId" binding:"required,min=1" label:"部署脚本"`
	HostGroupId int64  `json:"hostGroupId" binding:"required,min=1" label:"服务器组"`
}
