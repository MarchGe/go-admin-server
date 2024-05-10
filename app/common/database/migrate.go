package database

import (
	autoMigrate "github.com/MarchGe/go-admin-server/app/admin/model/auto_migrate"
	"gorm.io/gorm"
)

func TableAutoMigrate(db *gorm.DB) error {
	return autoMigrate.TableAutoMigrate(db)
}
