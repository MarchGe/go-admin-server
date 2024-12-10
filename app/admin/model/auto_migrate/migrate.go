package auto_migrate

import (
	"github.com/MarchGe/go-admin-server/app/admin/model"
	"github.com/MarchGe/go-admin-server/app/admin/model/dvmodel"
	"github.com/MarchGe/go-admin-server/app/admin/model/dvmodel/task"
	"gorm.io/gorm"
)

func TableAutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Dept{},
		&model.Icon{},
		&model.Job{},
		&model.ExceptionLog{},
		&model.LoginLog{},
		&model.OpLog{},
		&model.Menu{},
		&model.Role{},
		&model.RoleMenu{},
		&model.RoleUser{},
		&model.Settings{},
		&model.User{},
		&model.UserJob{},
		&model.UserMenu{},
		&model.UserPassword{},
		&dvmodel.App{},
		&task.DeployTask{},
		&dvmodel.Group{},
		&dvmodel.Host{},
		&dvmodel.HostGroup{},
		&dvmodel.Script{},
		&task.Task{},
		&task.ScriptTask{},
		&task.ScriptTaskScript{},
	)
}
