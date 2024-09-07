package req

type SftpRenameReq struct {
	ExplorerRenameReq
	HostId int64 `json:"hostId" binding:"required,min=1" label:"主机主键ID"`
}
