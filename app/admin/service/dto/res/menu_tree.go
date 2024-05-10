package res

import (
	"github.com/MarchGe/go-admin-server/app/admin/model"
)

type MenuTree struct {
	model.Menu
	Children []*MenuTree `json:"children"`
}
