package dvroutes

import (
	"github.com/MarchGe/go-admin-server/app"
	"github.com/MarchGe/go-admin-server/app/admin/apis/devops/xterm"
	"github.com/MarchGe/go-admin-server/app/common/middleware/authz"
	"github.com/gin-gonic/gin"
)

func init() {
	app.RouterRegister(registerXtermRoutes)
}

func registerXtermRoutes(g *gin.RouterGroup) {
	a := xterm.GetXterm()
	rg := g.Group("/terminal")

	rg.GET("/ws", authz.RequiresPermissions("terminal:connect"), a.Connect)
	rg.GET("/ws/ssh/:id", authz.RequiresPermissions("host:connect"), a.ConnectWithRemoteSSH)
}
