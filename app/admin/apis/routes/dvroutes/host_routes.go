package dvroutes

import (
	"github.com/MarchGe/go-admin-server/app"
	"github.com/MarchGe/go-admin-server/app/admin/apis/devops"
	"github.com/MarchGe/go-admin-server/app/common/middleware/authz"
	"github.com/MarchGe/go-admin-server/app/common/middleware/recorder"
	"github.com/gin-gonic/gin"
)

func init() {
	app.RouterRegister(registerHostRoutes)
}

func registerHostRoutes(g *gin.RouterGroup) {
	a := devops.GetHostApi()
	rg := g.Group("/devops/host")

	rg.POST("", authz.RequiresPermissions("host:add"), recorder.RecordOpLog("服务器", true), a.AddHost)
	rg.PUT("/:id", authz.RequiresPermissions("host:update"), recorder.RecordOpLog("服务器", true), a.UpdateHost)
	rg.DELETE("/:id", authz.RequiresPermissions("host:delete"), recorder.RecordOpLog("服务器"), a.DeleteHost)
	rg.GET("/list", authz.RequiresPermissions("host:list"), a.GetList)
	rg.GET("/all", authz.RequiresPermissions("host:list", "group:list"), a.GetAll)
	rg.GET("/connect-test", authz.RequiresPermissions("host:connectTest"), a.SshConnectTest)
}
