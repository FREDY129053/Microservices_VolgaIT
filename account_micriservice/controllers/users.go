package controllers

import (
	"account_microservice/models"
	"account_microservice/database"
	"account_microservice/helpers"
	"github.com/gin-gonic/gin"
	"log"
)

var databaseConn = database.GetConnection()

func Signup(c *gin.Context) {
	var user models.SignupUser

	// Проверка на правильность переданных данных
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		c.Abort()
		return
	}

	// Проверка на существование пользователя
	var existingUser models.SigninUser
	row := databaseConn.QueryRow("SELECT username,password FROM Users WHERE username=$1", user.Username)
	if err := row.Scan(&existingUser.Username, &existingUser.Password); err == nil {
		c.JSON(400, gin.H{"message": "User already exists"})
		c.Abort()
		return
	}
	
	userUUID := helpers.GenerateUUID()
	_, err := databaseConn.Exec("INSERT INTO user_and_roles (user_uuid) VALUES ($1)", userUUID)
	if err != nil {
		log.Println(err.Error())
		c.JSON(500, gin.H{"message": "Cannot create user"})
		c.Abort()
		return
	}
	_, err2 := databaseConn.Exec("INSERT INTO users (uuid, username, first_name, last_name, password) VALUES($1, $2, $3, $4, $5)", userUUID, user.Username, user.FirstName, user.LastName, user.Password)
	if err2 != nil {
		log.Println(err2.Error())
		c.JSON(500, gin.H{"message": "Cannot create user"})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{"message": "User created successfully"})
}