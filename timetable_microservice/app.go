package main

import (
	"time"
	"timetable_microservice/controllers"
	"timetable_microservice/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	defer router.Run("127.0.0.1:8083")

	// CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://127.0.0.1:8081", "http://127.0.0.1:8082", "http://127.0.0.1:8083", "http://127.0.0.1:8084"}
	config.AllowMethods = []string{"POST", "GET", "PUT", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour
	router.Use(cors.New(config))

	hospitals := router.Group("/api")
	hospitals.POST("/Timetable", middlewares.IsAdminOrManager(), controllers.AddNewNote)
	hospitals.PUT("/Timetable/:id", middlewares.IsAdminOrManager(), controllers.UpdateNote)
	hospitals.DELETE("/Timetable/:id", middlewares.IsAdminOrManager(), controllers.DeleteByID)
	hospitals.DELETE("/Timetable/Doctor/:uuid", middlewares.IsAdminOrManager(), controllers.DeleteByDoctorID)
	hospitals.DELETE("/Timetable/Hospital/:uuid", middlewares.IsAdminOrManager(), controllers.DeleteByHospitalID)

	hospitals.GET("/Timetable/Hospital/:uuid", middlewares.IsAuthorized(), controllers.GetByHospitalUUID)
	hospitals.GET("/Timetable/Doctor/:uuid", middlewares.IsAuthorized(), controllers.GetByDoctorUUID)
	hospitals.GET("/Timetable/Hospital/:uuid/Room/:room", middlewares.IsAdminOrManagerOrDoctor(), controllers.GetByHospitalUUIDAndRoom)
	hospitals.GET("/Timetable/:id/Appointments", middlewares.IsAuthorized(), controllers.GetAppointments)
	hospitals.POST("/Timetable/:id/Appointments", middlewares.IsAuthorized(), controllers.MakeAnAppointment)
	hospitals.DELETE("/Appointment/:id", middlewares.IsAdminOrManagerOrPatient(), controllers.DeleteAppointment)
}