package res

import (
	"github.com/MarchGe/go-admin-server/app/admin/model"
)

type PageableUsers struct {
	List  []*model.User `json:"list"`
	Total int64         `json:"total"`
}
