package main

import (
	"hospital_microservice/middlewares"
	"hospital_microservice/controllers"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	defer router.Run("127.0.0.1:8082")

	// CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://127.0.0.1:8081", "http://127.0.0.1:8082"}
	config.AllowMethods = []string{"POST", "GET", "PUT", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour
	router.Use(cors.New(config))

	hospitals := router.Group("/api/Hospitals")
	hospitals.GET("/", middlewares.IsAuthorized(), controllers.GetAllHospitals)
	hospitals.GET("/:uuid", middlewares.IsAuthorized(), controllers.GetHospitalInfo)
	hospitals.GET("/:uuid/Rooms", middlewares.IsAuthorized(), controllers.GetHospitalRooms)
	hospitals.POST("/", middlewares.IsAdmin(), controllers.AddHospital)
	hospitals.PUT("/:uuid", middlewares.IsAdmin(), controllers.UpdateHospital)
	hospitals.DELETE("/:uuid", middlewares.IsAdmin(), controllers.DeleteHospital)
}