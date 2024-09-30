package main

import (
	"document_microservice/controllers"
	"document_microservice/middlewares"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	defer router.Run("127.0.0.1:8084")

	// CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://127.0.0.1:8081", "http://127.0.0.1:8082", "http://127.0.0.1:8083", "http://127.0.0.1:8084"}
	config.AllowMethods = []string{"POST", "GET", "PUT", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour
	router.Use(cors.New(config))

	history := router.Group("/api/History")
	history.GET("/Account/:uuid", middlewares.IsDoctorOrPatient(), controllers.GetAllAccountHistories)
	history.GET("/:id", middlewares.IsDoctorOrPatient(), controllers.GetHistory)

	history.POST("/", middlewares.IsAdminOrManagerOrDoctor(), controllers.AddNewHistory)
	history.PUT("/:id", middlewares.IsAdminOrManagerOrDoctor(), controllers.UpdateHistory)
}