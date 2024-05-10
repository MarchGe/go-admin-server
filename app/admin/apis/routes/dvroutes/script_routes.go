package dvroutes

import (
	"github.com/MarchGe/go-admin-server/app"
	"github.com/MarchGe/go-admin-server/app/admin/apis/devops"
	"github.com/MarchGe/go-admin-server/app/common/middleware/authz"
	"github.com/MarchGe/go-admin-server/app/common/middleware/recorder"
	"github.com/gin-gonic/gin"
)

func init() {
	app.RouterRegister(registerScriptRoutes)
}

func registerScriptRoutes(g *gin.RouterGroup) {
	a := devops.GetScriptApi()
	rg := g.Group("/devops/script")

	rg.POST("", authz.RequiresPermissions("script:add"), recorder.RecordOpLog("脚本"), a.AddScript)
	rg.PUT("/:id", authz.RequiresPermissions("script:update"), recorder.RecordOpLog("脚本"), a.UpdateScript)
	rg.DELETE("/:id", authz.RequiresPermissions("script:delete"), recorder.RecordOpLog("脚本"), a.DeleteScript)
	rg.GET("/list", authz.RequiresPermissions("script:list", "task:list"), a.GetList)
}
