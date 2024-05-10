package routes

import (
	"github.com/MarchGe/go-admin-server/app"
	"github.com/MarchGe/go-admin-server/app/common/constant"
	"github.com/MarchGe/go-admin-server/config"
	"github.com/MarchGe/go-admin-server/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	app.RouterRegister(registerSwaggerRoutes)
}

func registerSwaggerRoutes(g *gin.RouterGroup) {
	rg := g.Group(constant.Swagger)

	docs.SwaggerInfo.BasePath = config.GetConfig().ContextPath
	rg.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
