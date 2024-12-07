package res

import (
	"github.com/MarchGe/go-admin-server/app/admin/model/dvmodel/task"
)

type TaskRes struct {
	*task.Task
	Concrete any `json:"concrete"`
}
