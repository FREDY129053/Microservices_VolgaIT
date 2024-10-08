package main

import (
	"time"

	"account_microservice/controllers"
	"account_microservice/middlewares"
	_ "account_microservice/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/files"
)

// @title Account microservice API
// @version 1.0
// @description Account API on Go documentation

// @host 0.0.0.0:8081
// @BasePath /api
func main() {
	router := gin.Default()
	defer router.Run("0.0.0.0:8081")

	// CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://127.0.0.1:8081", "http://127.0.0.1:8082", "http://127.0.0.1:8083", "http://127.0.0.1:8084",
		"http://0.0.0.0:8081", "http://0.0.0.0:8082", "http://0.0.0.0:8083", "http://0.0.0.0:8084",
	}
	config.AllowMethods = []string{"POST", "GET", "PUT", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour
	router.Use(cors.New(config))

	accounts := router.Group("/api")
	accounts.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	accounts.POST("/Authentication/SignUp", controllers.Signup) // +
	accounts.POST("/Authentication/SignIn", controllers.Signin) // +
	accounts.PUT("/Authentication/SignOut", middlewares.IsAuthorized(), controllers.SignOut)
	accounts.GET("/Authentication/Validate", controllers.VerifyingToken) // +
	accounts.PUT("/Authentication/Refresh", controllers.RefreshAccessToken) // +
	accounts.GET("/Accounts/Me", middlewares.IsAuthorized(), controllers.GetInfoAboutAccount) // +
	accounts.PUT("/Accounts/Update", middlewares.IsAuthorized(), controllers.UpdateAccount) // +
	accounts.GET("/Accounts", middlewares.IsAdmin(), controllers.GetAccounts) // +
	accounts.POST("/Accounts", middlewares.IsAdmin(), controllers.AddAccountByAdmin) // +
	accounts.PUT("/Accounts/:uuid", middlewares.IsAdmin(), controllers.ChangeAccountByAdmin) // +
	accounts.DELETE("/Accounts/:uuid", middlewares.IsAdmin(), controllers.DeleteAccountByAdmin) // +
	accounts.GET("/Accounts/Doctors", middlewares.IsAuthorized(), controllers.GetAllDoctors) // +
	accounts.GET("/Accounts/Doctors/:uuid", middlewares.IsAuthorized(), controllers.GetDoctor)
}
