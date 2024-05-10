package res

import (
	"github.com/MarchGe/go-admin-server/app/admin/model"
)

type DeptTree struct {
	model.Dept
	Children []*DeptTree `json:"children"`
}
