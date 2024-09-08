package routes

import (
	"example.com/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	//
	server.POST("/signup", signup)
	server.POST("/signin", signin)
	//
	server.POST("/refresh", middlewares.Refresh, refresh)
	//
	authenticated := server.Group("/")
	authenticated.Use(middlewares.Authenticate)
	//
}
