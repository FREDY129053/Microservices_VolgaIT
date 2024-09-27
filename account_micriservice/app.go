package main

import (
	"time"

	"account_microservice/controllers"
	"account_microservice/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	defer router.Run("127.0.0.1:8081")

	// CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"POST", "GET", "PUT", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour
	router.Use(cors.New(config))

	accounts := router.Group("/api")
	accounts.POST("/Authentication/SignUp", controllers.Signup)
	accounts.POST("/Authentication/SignIn", controllers.Signin)
	accounts.PUT("/Authentication/SignOut", middlewares.IsAuthorized(), controllers.SignOut)
	accounts.GET("/Authentication/Validate", controllers.VerifyingToken)
	accounts.POST("/Authentication/Refresh", controllers.RefreshAccessToken)
	accounts.GET("/Accounts/Me", middlewares.IsAuthorized(), controllers.GetInfoAboutAccount)
	accounts.PUT("/Accounts/Update", middlewares.IsAuthorized(), controllers.UpdateAccount)
	accounts.GET("/Accounts", middlewares.IsAdmin(), controllers.GetAccounts)
	accounts.POST("/Accounts", middlewares.IsAdmin(), controllers.AddAccountByAdmin)
	accounts.PUT("/Accounts/:uuid", middlewares.IsAdmin(), controllers.ChangeAccountByAdmin)
	accounts.DELETE("/Accounts/:uuid", middlewares.IsAdmin(), controllers.DeleteAccountByAdmin)
	accounts.GET("/Accounts/Doctors", middlewares.IsAuthorized(), controllers.GetAllDoctors)
	accounts.GET("/Accounts/Doctors/:uuid", middlewares.IsAuthorized(), controllers.GetDoctor)
}
