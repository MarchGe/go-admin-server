package res

type UploadRes struct {
	TmpKey   string `json:"tmpKey"`
	FileName string `json:"fileName"`
	FileSize int64  `json:"fileSize"`
}
