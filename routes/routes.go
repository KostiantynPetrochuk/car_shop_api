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
	server.GET("/brands", getBrands)
	server.POST("/brands", addBrand)
	//
	server.GET("/features", getFeatures)
	server.POST("/features", addFeature)
	//
	server.POST("/models", addModel)
	//
	authenticated := server.Group("/")
	authenticated.Use(middlewares.Authenticate)
	//
	authenticated.GET("/cars", getCars)
	authenticated.GET("/cars/:id", getCar)
	authenticated.POST("/cars", addCar)
}
