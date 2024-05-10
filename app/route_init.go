package app

import (
	"github.com/gin-gonic/gin"
)

var fns = make([]func(g *gin.RouterGroup), 0)

type RouterRegisterFunc func(g *gin.RouterGroup)

// InitRoutes 初始化路由信息
func InitRoutes(g *gin.RouterGroup) {
	for _, fn := range fns {
		fn(g)
	}
}

// RouterRegister 注册路由信息
func RouterRegister(f RouterRegisterFunc) {
	fns = append(fns, f)
}
