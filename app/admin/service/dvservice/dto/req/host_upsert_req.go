package req

type HostUpsertReq struct {
	Name            string `json:"name" binding:"required,max=50" label:"名称"`
	Ip              string `json:"ip" binding:"required,max=50" label:"ip地址"`
	Port            int32  `json:"port" binding:"required,min=1,max=65535" label:"端口"`
	User            string `json:"user" binding:"required,max=50" label:"用户名"`
	Password        string `json:"password" binding:"required" label:"密码"`
	SortNum         int32  `json:"sortNum" binding:"omitempty" label:"排序"`
	PasswordChanged bool   `json:"passwordChanged" binding:"omitempty"` // 用于标识密码是否被手动修改
}
