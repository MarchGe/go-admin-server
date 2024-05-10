package req

const AppPkgFileNameMaxLength = 100 // 跟结构体里的FileName字段长度限制保持一致
type AppUpsertReq struct {
	Name     string `json:"name" binding:"required,max=50" label:"名称"`
	Version  string `json:"version" binding:"required,max=50" label:"版本"`
	Port     int16  `json:"port" binding:"required,min=1,max=65535" label:"端口"`
	Key      string `json:"key" binding:"required,max=255" label:"部署包"`
	FileName string `json:"fileName" binding:"required,max=100" label:"文件名"`
}
