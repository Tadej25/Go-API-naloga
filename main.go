package main

// @title Go API naloga
// @version 1.0
// @description Basic API for displaying temperature data for cities.
// @BasePath /
// @Security basicAuth

import (
	"fmt"

	_ "github.com/ahmetb/go-linq/v3"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"

	"exmaple/Go-API-naloga/controllers"
	_ "exmaple/Go-API-naloga/docs"

	ginSwagger "github.com/swaggo/gin-swagger"
)

// @securityDefinitions.basic BasicAuth
// @in header
// @name Authorization
func main() {
	err := controllers.ReadCsv()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return
	}
	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/cities", controllers.GetCities)
	router.GET("/city/:name", controllers.GetCityByName)
	router.GET("/AverageTemperatures", controllers.GetAverageTemperatures)
	router.POST("reload", controllers.AuthMiddleware(), controllers.Reload)
	router.Run("localhost:8080")
}
