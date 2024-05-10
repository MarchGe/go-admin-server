package req

const HostUpdateMode = 1

type SshConnectTestParams struct {
	Ip              string `form:"ip" json:"ip" binding:"required,max=50" label:"ip地址"`
	Port            int16  `form:"port" json:"port" binding:"required,min=1,max=65535" label:"端口"`
	User            string `form:"user" json:"user" binding:"required,max=50" label:"用户名"`
	Password        string `form:"password" json:"password" binding:"required" label:"密码"`
	Mode            int8   `form:"mode" json:"mode" binding:"omitempty,oneof=0 1" label:"模式"`
	PasswordChanged bool   `form:"passwordChanged" json:"passwordChanged" binding:"omitempty"` // 用于标识密码是否被手动修改
}
