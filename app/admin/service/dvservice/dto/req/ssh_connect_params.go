package req

type SshConnectParams struct {
	Ip       string `form:"ip" json:"ip" binding:"required,max=50" label:"ip地址"`
	Port     int32  `form:"port" json:"port" binding:"required,min=1,max=65535" label:"端口"`
	User     string `form:"user" json:"user" binding:"required,max=50" label:"用户名"`
	Password string `form:"password" json:"password" binding:"required" label:"密码"`
}
