package routes

import (
	"github.com/MarchGe/go-admin-server/app"
	"github.com/MarchGe/go-admin-server/app/admin/apis"
	"github.com/MarchGe/go-admin-server/app/common/middleware/authz"
	"github.com/gin-gonic/gin"
)

func init() {
	app.RouterRegister(registerIconRoutes)
}

func registerIconRoutes(g *gin.RouterGroup) {
	a := apis.GetIconApi()
	rg := g.Group("/icon")

	rg.GET("/all", authz.RequiresPermissions("icon:all", "menu:tree"), a.GetAll)
}
