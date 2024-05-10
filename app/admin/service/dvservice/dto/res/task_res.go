package res

import "github.com/MarchGe/go-admin-server/app/admin/model/dvmodel"

type TaskRes struct {
	*dvmodel.Task
	Concrete any `json:"concrete"`
}
