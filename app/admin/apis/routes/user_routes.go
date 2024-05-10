package routes

import (
	"github.com/MarchGe/go-admin-server/app"
	"github.com/MarchGe/go-admin-server/app/admin/apis"
	"github.com/MarchGe/go-admin-server/app/common/middleware/authz"
	"github.com/MarchGe/go-admin-server/app/common/middleware/recorder"
	"github.com/gin-gonic/gin"
)

func init() {
	app.RouterRegister(registerUserRoutes)
}

func registerUserRoutes(g *gin.RouterGroup) {
	a := apis.GetUserApi()
	rg := g.Group("/user")

	rg.POST("", authz.RequiresPermissions("user:add"), recorder.RecordOpLog("用户"), a.AddUser)
	rg.PUT("/:id", authz.RequiresPermissions("user:update"), recorder.RecordOpLog("用户"), a.UpdateUser)
	rg.DELETE("/:id", authz.RequiresPermissions("user:delete"), recorder.RecordOpLog("用户"), a.DeleteUser)
	rg.GET("/:id", authz.RequiresPermissions("user:get"), a.GetUser)
	rg.GET("/list", authz.RequiresPermissions("user:list"), a.GetList)
	rg.PUT("/enable/:id", authz.RequiresPermissions("user:enable"), recorder.RecordOpLog("用户", "启用"), a.EnableAccount)
	rg.PUT("/disable/:id", authz.RequiresPermissions("user:disable"), recorder.RecordOpLog("用户", "禁用"), a.DisableAccount)
	rg.PUT("/passwd", authz.RequiresPermissions("user:changePassword"), recorder.RecordOpLog("用户密码", true), a.ChangePasswd)
	rg.PUT("/passwd-reset/:id", authz.RequiresPermissions("user:resetPassword"), recorder.RecordOpLog("用户密码", "重置"), a.ResetPasswd)
	rg.PUT("/:id/menus", authz.RequiresPermissions("user:menus"), recorder.RecordOpLog("用户权限"), a.UpdateUserMenus)
}
