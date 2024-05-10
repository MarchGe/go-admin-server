package req

type SettingsUpsertReq struct {
	Key   string `json:"key" binding:"required,max=50" label:"键名"`
	Value string `json:"value" binding:"required,max=2048" label:"值"`
}
