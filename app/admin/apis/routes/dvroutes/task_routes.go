package dvroutes

import (
	"github.com/MarchGe/go-admin-server/app"
	"github.com/MarchGe/go-admin-server/app/admin/apis/devops"
	"github.com/MarchGe/go-admin-server/app/common/middleware/authz"
	"github.com/MarchGe/go-admin-server/app/common/middleware/recorder"
	"github.com/gin-gonic/gin"
)

func init() {
	app.RouterRegister(registerTaskRoutes)
}

func registerTaskRoutes(g *gin.RouterGroup) {
	a := devops.GetTaskApi()
	rg := g.Group("/devops/task")

	rg.POST("", authz.RequiresPermissions("task:add"), recorder.RecordOpLog("任务"), a.AddTask)
	rg.PUT("/:id", authz.RequiresPermissions("task:update"), recorder.RecordOpLog("任务"), a.UpdateTask)
	rg.DELETE("/:id", authz.RequiresPermissions("task:delete"), recorder.RecordOpLog("任务"), a.DeleteTask)
	rg.GET("/list", authz.RequiresPermissions("task:list"), a.GetList)
	rg.POST("/start/:id", authz.RequiresPermissions("task:start"), recorder.RecordOpLog("任务", "启动"), a.StartTask)
	rg.POST("/stop/:id", authz.RequiresPermissions("task:stop"), recorder.RecordOpLog("任务", "停止"), a.StopTask)
}
