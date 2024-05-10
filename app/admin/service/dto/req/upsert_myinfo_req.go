package req

type UpsertMyInfoReq struct {
	Cellphone string `json:"cellphone" binding:"omitempty,regex=^1[0-9]{10}$" label:"手机号"`
	Nickname  string `json:"nickname" binding:"required,max=20" label:"昵称"`
}
