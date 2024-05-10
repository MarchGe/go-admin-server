package routes

import (
	"github.com/MarchGe/go-admin-server/app"
	"github.com/MarchGe/go-admin-server/app/admin/apis"
	"github.com/MarchGe/go-admin-server/app/common/middleware/authz"
	"github.com/MarchGe/go-admin-server/app/common/middleware/recorder"
	"github.com/gin-gonic/gin"
)

func init() {
	app.RouterRegister(registerMenuRoutes)
}

func registerMenuRoutes(g *gin.RouterGroup) {
	a := apis.GetMenuApi()
	rg := g.Group("/menu")

	rg.POST("", authz.RequiresPermissions("menu:add"), recorder.RecordOpLog("菜单"), a.AddMenu)
	rg.PUT("/:id", authz.RequiresPermissions("menu:update"), recorder.RecordOpLog("菜单"), a.UpdateMenu)
	rg.DELETE("/:id", authz.RequiresPermissions("menu:delete"), recorder.RecordOpLog("菜单"), a.DeleteMenu)
	rg.GET("/:id", authz.RequiresPermissions("menu:get"), a.GetMenu)
	rg.GET("/tree", authz.RequiresPermissions("menu:tree", "user:list", "role:list"), a.GetMenuTree)
	rg.GET("/my/tree", a.GetMyMenuTree)
	rg.GET("/my/permissions", a.GetMyPermissions)
}
