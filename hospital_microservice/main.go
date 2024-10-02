package main

import (
	"hospital_microservice/middlewares"
	"hospital_microservice/controllers"
	_ "hospital_microservice/docs"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/files"
)

// @title Hospital microservice API
// @version 1.0
// @description Hospital API on Go documentation

// @host localhost:8082
// @BasePath /api/Hospitals
func main() {
	router := gin.Default()
	defer router.Run("127.0.0.1:8082")

	// CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://127.0.0.1:8081", "http://127.0.0.1:8082", "http://127.0.0.1:8083", "http://127.0.0.1:8084"}
	config.AllowMethods = []string{"POST", "GET", "PUT", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour
	router.Use(cors.New(config))

	hospitals := router.Group("/api/Hospitals")
	hospitals.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	hospitals.GET("/", middlewares.IsAuthorized(), controllers.GetAllHospitals)
	hospitals.GET("/:uuid", middlewares.IsAuthorized(), controllers.GetHospitalInfo)
	hospitals.GET("/:uuid/Rooms", middlewares.IsAuthorized(), controllers.GetHospitalRooms)
	hospitals.POST("/", middlewares.IsAdmin(), controllers.AddHospital)
	hospitals.PUT("/:uuid", middlewares.IsAdmin(), controllers.UpdateHospital)
	hospitals.DELETE("/:uuid", middlewares.IsAdmin(), controllers.DeleteHospital)
}