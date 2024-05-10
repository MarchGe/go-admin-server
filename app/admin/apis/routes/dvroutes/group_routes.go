package dvroutes

import (
	"github.com/MarchGe/go-admin-server/app"
	"github.com/MarchGe/go-admin-server/app/admin/apis/devops"
	"github.com/MarchGe/go-admin-server/app/common/middleware/authz"
	"github.com/MarchGe/go-admin-server/app/common/middleware/recorder"
	"github.com/gin-gonic/gin"
)

func init() {
	app.RouterRegister(registerGroupRoutes)
}

func registerGroupRoutes(g *gin.RouterGroup) {
	a := devops.GetGroupApi()
	rg := g.Group("/devops/group")

	rg.POST("", authz.RequiresPermissions("group:add"), recorder.RecordOpLog("服务器组"), a.AddGroup)
	rg.PUT("/:id", authz.RequiresPermissions("group:update"), recorder.RecordOpLog("服务器组"), a.UpdateGroup)
	rg.DELETE("/:id", authz.RequiresPermissions("group:delete"), recorder.RecordOpLog("服务器组"), a.DeleteGroup)
	rg.GET("/list", authz.RequiresPermissions("group:list", "task:list"), a.GetList)
}
