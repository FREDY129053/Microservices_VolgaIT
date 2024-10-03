package main

import (
	"time"
	"timetable_microservice/controllers"
	_ "timetable_microservice/docs"
	"timetable_microservice/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Timetable microservice API
// @version 1.0
// @description Timetable API on Go documentation

// @host 0.0.0.0:8083
// @BasePath /api
func main() {
	router := gin.Default()
	defer router.Run("0.0.0.0:8083")

	// CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://127.0.0.1:8081", "http://127.0.0.1:8082", "http://127.0.0.1:8083", "http://127.0.0.1:8084", "http://0.0.0.0:8081", "http://0.0.0.0:8082", "http://0.0.0.0:8083", "http://0.0.0.0:8084",}
	config.AllowMethods = []string{"POST", "GET", "PUT", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour
	router.Use(cors.New(config))

	timetable := router.Group("/api")
	timetable.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	timetable.POST("/Timetable", middlewares.IsAdminOrManager(), controllers.AddNewNote)
	timetable.PUT("/Timetable/:id", middlewares.IsAdminOrManager(), controllers.UpdateNote)
	timetable.DELETE("/Timetable/:id", middlewares.IsAdminOrManager(), controllers.DeleteByID)
	timetable.DELETE("/Timetable/Doctor/:uuid", middlewares.IsAdminOrManager(), controllers.DeleteByDoctorID)
	timetable.DELETE("/Timetable/Hospital/:uuid", middlewares.IsAdminOrManager(), controllers.DeleteByHospitalID)

	timetable.GET("/Timetable/Hospital/:uuid", middlewares.IsAuthorized(), controllers.GetByHospitalUUID)
	timetable.GET("/Timetable/Doctor/:uuid", middlewares.IsAuthorized(), controllers.GetByDoctorUUID)
	timetable.GET("/Timetable/Hospital/:uuid/Room/:room", middlewares.IsAdminOrManagerOrDoctor(), controllers.GetByHospitalUUIDAndRoom)
	timetable.GET("/Timetable/:id/Appointments", middlewares.IsAuthorized(), controllers.GetAppointments)
	timetable.POST("/Timetable/:id/Appointments", middlewares.IsAuthorized(), controllers.MakeAnAppointment)
	timetable.DELETE("/Appointment/:id", middlewares.IsAdminOrManagerOrPatient(), controllers.DeleteAppointment)
}