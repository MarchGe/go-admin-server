package routes

import (
	"github.com/MarchGe/go-admin-server/app"
	"github.com/MarchGe/go-admin-server/app/admin/apis"
	"github.com/gin-gonic/gin"
)

func init() {
	app.RouterRegister(registerSettingsRoutes)
}

func registerSettingsRoutes(g *gin.RouterGroup) {
	a := apis.GetSettingsApi()
	rg := g.Group("/settings")

	rg.GET("", a.GetSettings)
	rg.PUT("", a.Upsert)
}
