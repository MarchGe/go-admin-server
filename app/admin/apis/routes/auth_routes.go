package routes

import (
	"github.com/MarchGe/go-admin-server/app"
	"github.com/MarchGe/go-admin-server/app/admin/apis"
	"github.com/MarchGe/go-admin-server/app/common/middleware/authz"
	"github.com/gin-gonic/gin"
)

func init() {
	app.RouterRegister(registerAuthRoutes)
}

func registerAuthRoutes(g *gin.RouterGroup) {
	a := apis.GetAuthApi()
	rg := g.Group("/auth")

	rg.POST("/login", a.Login)
	rg.PUT("/passwd", a.ChangeMyPasswd)
	rg.GET("/my-info", a.GetMyInfo)
	rg.PUT("/my-info", a.UpdateMyInfo)
	rg.GET("/web-shell-token", authz.RequiresPermissions("terminal:connect"), a.GetWebShellToken)
}
