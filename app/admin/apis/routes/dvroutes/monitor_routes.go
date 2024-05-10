package dvroutes

import (
	"github.com/MarchGe/go-admin-server/app"
	"github.com/MarchGe/go-admin-server/app/admin/apis/devops/monitor"
	"github.com/MarchGe/go-admin-server/app/common/middleware/authz"
	"github.com/gin-gonic/gin"
)

func init() {
	app.RouterRegister(registerMonitorRoutes)
}

func registerMonitorRoutes(g *gin.RouterGroup) {
	a := monitor.GetSysStatsApi()
	rg := g.Group("/monitor")

	rg.GET("/list", authz.RequiresPermissions("monitor:list"), a.GetList)
	rg.GET("/performance-stats", authz.RequiresPermissions("monitor:detail"), a.GetPerformanceStats)
	rg.GET("/host-info", authz.RequiresPermissions("monitor:detail"), a.GetHostInfo)
	rg.DELETE("/host", authz.RequiresPermissions("monitor:delete"), a.DeleteHost)
}
