package routes

import (
	"github.com/MarchGe/go-admin-server/app"
	"github.com/MarchGe/go-admin-server/app/admin/apis"
	"github.com/MarchGe/go-admin-server/app/common/middleware/authz"
	"github.com/MarchGe/go-admin-server/app/common/middleware/recorder"
	"github.com/gin-gonic/gin"
)

func init() {
	app.RouterRegister(registerDeptRoutes)
}

func registerDeptRoutes(g *gin.RouterGroup) {
	a := apis.GetDeptApi()
	rg := g.Group("/dept")

	rg.POST("", authz.RequiresPermissions("dept:add"), recorder.RecordOpLog("部门"), a.AddDept)
	rg.PUT("/:id", authz.RequiresPermissions("dept:update"), recorder.RecordOpLog("部门"), a.UpdateDept)
	rg.DELETE("/:id", authz.RequiresPermissions("dept:delete"), recorder.RecordOpLog("部门"), a.DeleteDept)
	rg.GET("/:id", authz.RequiresPermissions("dept:get"), a.GetDept)
	rg.GET("/tree", authz.RequiresPermissions("dept:tree", "user:list"), a.GetDeptTree)
}
