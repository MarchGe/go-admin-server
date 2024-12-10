package dvroutes

import (
	"github.com/MarchGe/go-admin-server/app"
	"github.com/MarchGe/go-admin-server/app/admin/apis/devops"
	"github.com/MarchGe/go-admin-server/app/common/middleware/authz"
	"github.com/MarchGe/go-admin-server/app/common/middleware/recorder"
	"github.com/gin-gonic/gin"
)

func init() {
	app.RouterRegister(registerScripTaskRoutes)
}

func registerScripTaskRoutes(g *gin.RouterGroup) {
	a := devops.GetScriptTaskApi()
	rg := g.Group("/devops/script-task")

	rg.POST("", authz.RequiresPermissions("scriptTask:add"), recorder.RecordOpLog("脚本任务"), a.Add)
	rg.PUT("/:id", authz.RequiresPermissions("scriptTask:update"), recorder.RecordOpLog("脚本任务"), a.Update)
	rg.DELETE("/:id", authz.RequiresPermissions("scriptTask:delete"), recorder.RecordOpLog("脚本任务"), a.Delete)
	rg.GET("/list", authz.RequiresPermissions("scriptTask:list"), a.GetList)
	rg.POST("/start/:id", authz.RequiresPermissions("scriptTask:start"), recorder.RecordOpLog("脚本任务", "启动"), a.Start)
	rg.POST("/stop/:id", authz.RequiresPermissions("scriptTask:stop"), recorder.RecordOpLog("脚本任务", "停止"), a.Stop)
}
