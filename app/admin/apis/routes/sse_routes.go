package routes

import (
	"github.com/MarchGe/go-admin-server/app"
	"github.com/MarchGe/go-admin-server/app/admin/apis"
	"github.com/MarchGe/go-admin-server/app/common/middleware/authz"
	"github.com/gin-gonic/gin"
)

func init() {
	app.RouterRegister(registerSseRoutes)
}

func registerSseRoutes(g *gin.RouterGroup) {
	a := apis.GetSseApi()
	rg := g.Group("/sse")

	rg.GET("/message-push", a.MessagePush)
	rg.GET("/task/:id/manifest-log", authz.RequiresPermissions("task:log"), a.PushManifestLogEvent)
	rg.GET("/task/:id/host-log", authz.RequiresPermissions("task:log"), a.PushHostLogEvent)
	rg.GET("/script-task/:id/manifest-log", authz.RequiresPermissions("scriptTask:log"), a.PushScriptTaskManifestLog)
	rg.GET("/script-task/:id/host-log", authz.RequiresPermissions("scriptTask:log"), a.PushScriptTaskHostLogEvent)
}
