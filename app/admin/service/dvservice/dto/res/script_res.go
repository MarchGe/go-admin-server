package res

import "github.com/MarchGe/go-admin-server/app/admin/model/dvmodel"

type ScriptRes struct {
	dvmodel.Script
	DeployTaskRefCount int32 `json:"deployTaskRefCount"`
	ScriptTaskRefCount int32 `json:"scriptTaskRefCount"`
}
