package dvroutes

import (
	"github.com/MarchGe/go-admin-server/app"
	"github.com/MarchGe/go-admin-server/app/admin/apis/devops"
	"github.com/MarchGe/go-admin-server/app/common/middleware/authz"
	"github.com/MarchGe/go-admin-server/app/common/middleware/recorder"
	"github.com/gin-gonic/gin"
)

func init() {
	app.RouterRegister(registerExplorerRoutes)
}

func registerExplorerRoutes(g *gin.RouterGroup) {
	a := devops.GetExplorerApi()
	rg := g.Group("/devops/explorer")

	rg.DELETE("/entry", authz.RequiresPermissions("explorer:delete"), recorder.RecordOpLog("资源管理器资源"), a.DeleteEntry)
	rg.GET("/entries", authz.RequiresPermissions("explorer:entries"), a.GetEntries)
	rg.POST("/upload", authz.RequiresPermissions("explorer:upload"), a.Upload)
	rg.GET("/download", authz.RequiresPermissions("explorer:download"), a.Download)
}
