package req

type ChangePasswdReq struct {
	Id          int64  `json:"id" binding:"required,min=1" label:"用户ID"`
	NewPassword string `json:"newPassword" binding:"required,min=6,max=30" label:"新密码"`
}
