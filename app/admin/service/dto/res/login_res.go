package res

import (
	"github.com/MarchGe/go-admin-server/app/admin/model"
)

type LoginRes struct {
	model.User
	Password string `json:"password"`
}
