package routes

import (
	"github.com/MarchGe/go-admin-server/app"
	"github.com/MarchGe/go-admin-server/app/admin/apis"
	"github.com/MarchGe/go-admin-server/app/common/middleware/authz"
	"github.com/MarchGe/go-admin-server/app/common/middleware/recorder"
	"github.com/gin-gonic/gin"
)

func init() {
	app.RouterRegister(registerRoleRoutes)
}

func registerRoleRoutes(g *gin.RouterGroup) {
	a := apis.GetRoleApi()
	rg := g.Group("/role")

	rg.POST("", authz.RequiresPermissions("role:add"), recorder.RecordOpLog("角色"), a.AddRole)
	rg.PUT("/:id", authz.RequiresPermissions("role:update"), recorder.RecordOpLog("角色"), a.UpdateRole)
	rg.DELETE("/:id", authz.RequiresPermissions("role:delete"), recorder.RecordOpLog("角色"), a.DeleteRole)
	rg.GET("/:id", authz.RequiresPermissions("role:get"), a.GetRole)
	rg.GET("/list", authz.RequiresPermissions("role:list", "user:list"), a.GetList)
	rg.PUT("/:id/menus", authz.RequiresPermissions("role:menus"), recorder.RecordOpLog("角色权限"), a.UpdateRoleMenus)
}
