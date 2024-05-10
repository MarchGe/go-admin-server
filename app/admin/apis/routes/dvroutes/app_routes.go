package dvroutes

import (
	"github.com/MarchGe/go-admin-server/app"
	"github.com/MarchGe/go-admin-server/app/admin/apis/devops"
	"github.com/MarchGe/go-admin-server/app/common/middleware/authz"
	"github.com/MarchGe/go-admin-server/app/common/middleware/recorder"
	"github.com/gin-gonic/gin"
)

func init() {
	app.RouterRegister(registerAppRoutes)
}

func registerAppRoutes(g *gin.RouterGroup) {
	a := devops.GetAppApi()
	rg := g.Group("/devops/app")

	rg.POST("", authz.RequiresPermissions("app:add"), recorder.RecordOpLog("应用"), a.AddApp)
	rg.PUT("/:id", authz.RequiresPermissions("app:update"), recorder.RecordOpLog("应用"), a.UpdateApp)
	rg.DELETE("/:id", authz.RequiresPermissions("app:delete"), recorder.RecordOpLog("应用"), a.DeleteApp)
	rg.GET("/list", authz.RequiresPermissions("app:list", "task:list"), a.GetList)
	rg.POST("/upload", authz.RequiresPermissions("app:upload"), a.UploadPkg)
	rg.GET("/download", authz.RequiresPermissions("app:download"), a.DownloadPkg)
}
