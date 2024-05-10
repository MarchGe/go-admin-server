package routes

import (
	"github.com/MarchGe/go-admin-server/app"
	"github.com/MarchGe/go-admin-server/app/admin/apis"
	"github.com/MarchGe/go-admin-server/app/common/middleware/authz"
	"github.com/MarchGe/go-admin-server/app/common/middleware/recorder"
	"github.com/gin-gonic/gin"
)

func init() {
	app.RouterRegister(registerJobRoutes)
}

func registerJobRoutes(g *gin.RouterGroup) {
	a := apis.GetJobApi()
	rg := g.Group("/job")

	rg.POST("", authz.RequiresPermissions("job:add"), recorder.RecordOpLog("岗位"), a.AddJob)
	rg.PUT("/:id", authz.RequiresPermissions("job:update"), recorder.RecordOpLog("岗位"), a.UpdateJob)
	rg.DELETE("/:id", authz.RequiresPermissions("job:delete"), recorder.RecordOpLog("岗位"), a.DeleteJob)
	rg.GET("/list", authz.RequiresPermissions("job:list", "user:list"), a.GetList)
}
