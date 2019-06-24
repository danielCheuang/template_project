package router

import (
	"template_project/config"
	"template_project/handler"

	"github.com/gin-gonic/gin"
)

func InitRouters(e *gin.Engine, cfg config.Configuration) {
	rootRouterPrefix := cfg.Server.RootRouterPrefix
	if rootRouterPrefix == "" {
		rootRouterPrefix = "/api"
	}
	routerGroupAPI := e.Group(rootRouterPrefix)

	v1 := routerGroupAPI.Group("/v1")
	{
		v1.GET("/ping", handler.Ping)
		v1.POST("/test_post", handler.TestPost)
		v1.GET("/get_user", handler.GetUser)
	}
}
