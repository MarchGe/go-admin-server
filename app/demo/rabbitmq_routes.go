package demo

import (
	"github.com/MarchGe/go-admin-server/app"
	"github.com/MarchGe/go-admin-server/app/common/middleware/authz"
	"github.com/gin-gonic/gin"
)

func init() {
	app.RouterRegister(registerRabbitMQRoutes)
}

func registerRabbitMQRoutes(g *gin.RouterGroup) {
	rg := g.Group("/demo")

	a := GetRabbitApi()
	rg.POST("/rabbitmq/msg", authz.RequiresPermissions("demo:rabbitmq"), a.SendMsg)
}
