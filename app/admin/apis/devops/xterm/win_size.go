package xterm

type WinSize struct {
	Rows int `form:"rows" json:"rows" binding:"required,min=1" label:"行数"`
	Cols int `form:"cols" json:"cols" binding:"required,min=1" label:"列数"`
}
