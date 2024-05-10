package req

type LoginReq struct {
	Email    string `json:"email" binding:"required,email" label:"账号"`
	Password string `json:"password" binding:"required,max=30" label:"密码"`
}
