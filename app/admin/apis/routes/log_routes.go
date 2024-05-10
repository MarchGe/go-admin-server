package routes

import (
	"github.com/MarchGe/go-admin-server/app"
	"github.com/MarchGe/go-admin-server/app/admin/apis"
	"github.com/MarchGe/go-admin-server/app/common/middleware/authz"
	"github.com/gin-gonic/gin"
)

func init() {
	app.RouterRegister(registerLogRoutes)
}

func registerLogRoutes(g *gin.RouterGroup) {
	a := apis.GetLogApi()
	rg := g.Group("/log")

	rg.GET("/login", authz.RequiresPermissions("loginLog:list"), a.GetLoginLogList)
	rg.GET("/op", authz.RequiresPermissions("opLog:list"), a.GetOpLogList)
	rg.GET("/exception", authz.RequiresPermissions("exceptionLog:list"), a.GetExceptionLogList)
	rg.DELETE("/login", authz.RequiresPermissions("loginLog:delete"), a.DeleteLoginLog)
	rg.DELETE("/op", authz.RequiresPermissions("opLog:delete"), a.DeleteOpLog)
	rg.DELETE("/exception", authz.RequiresPermissions("exceptionLog:delete"), a.DeleteExceptionLog)
}
