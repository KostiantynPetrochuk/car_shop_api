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
	//
	server.GET("/features", getFeatures)
	//
	server.GET("/cars", getCars)
	server.GET("/cars/:id", getCar)
	server.GET("/cars/brand/:id", getCarsByBrand)
	server.GET("/featured-cars", getFeaturedCars)
	//
	//
	authenticated := server.Group("/")
	authenticated.Use(middlewares.Authenticate)
	//
	authenticated.POST("/cars", addCar)
	//
	authenticated.POST("/models", addModel)
	//
	authenticated.POST("/features", addFeature)
	//
	authenticated.POST("/brands", addBrand)
}
