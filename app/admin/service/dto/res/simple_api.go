package res

type SimpleApi struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

func CreateSimpleApi(method, path string) *SimpleApi {
	return &SimpleApi{
		Method: method,
		Path:   path,
	}
}
