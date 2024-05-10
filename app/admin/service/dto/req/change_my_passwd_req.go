package req

type ChangeMyPasswdReq struct {
	OldPassword string `json:"oldPassword" binding:"required,max=30" label:"旧密码"`
	NewPassword string `json:"newPassword" binding:"required,min=6,max=30" label:"新密码"`
}
