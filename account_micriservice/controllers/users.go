package controllers

import (
	"account_microservice/database"
	"account_microservice/helpers"
	"account_microservice/models"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var databaseConn = database.GetConnection()
var jwtKey = []byte(os.Getenv("SECRET_KEY"))

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

	// Вставка пользователя и его роли в БД
	userUUID := helpers.GenerateUUID()
	_, err := databaseConn.Exec("INSERT INTO user_and_roles (user_uuid) VALUES ($1)", userUUID)
	if err != nil {
		c.JSON(500, gin.H{"message": "Cannot create user"})
		c.Abort()
		return
	}
	_, err2 := databaseConn.Exec("INSERT INTO users (uuid, username, first_name, last_name, password) VALUES($1, $2, $3, $4, $5)", userUUID, user.Username, user.FirstName, user.LastName, user.Password)
	if err2 != nil {
		c.JSON(500, gin.H{"message": "Cannot create user"})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{"message": "User created successfully"})
}

func Signin(c *gin.Context) {
	var user models.SigninUser

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		c.Abort()
		return
	}

	// Проверка на наличие пользователя
	var existingUser models.SigninUser
	row := databaseConn.QueryRow("SELECT username,password FROM users WHERE username=$1", user.Username)
	if err := row.Scan(&existingUser.Username, &existingUser.Password); err != nil {
		c.JSON(400, gin.H{"message": "User does not exist"})
		c.Abort()
		return
	}

	// Проверка паролей при вводе и из БД
	if user.Password != existingUser.Password {
		c.JSON(400, gin.H{"message": "Invalid password"})
		c.Abort()
		return
	}

	// Получение UUID пользователя для последующего получения ролей
	var uuid string
	err := databaseConn.QueryRow("SELECT uuid FROM users WHERE username=$1", user.Username).Scan(&uuid)
	if err != nil {
		c.JSON(404, gin.H{"message": "User not found"})
		c.Abort()
		return
	}

	// Получаем роли
	var roles []string
	rolesRows, err := databaseConn.Query("SELECT role FROM user_and_roles WHERE user_uuid=$1", uuid)
	if err != nil {
		panic(err)
	}
	for rolesRows.Next() {
		var role string
		err := rolesRows.Scan(&role)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		roles = append(roles, role)
	}

	// Создание токена на 5 минут
	expirationTimeAccess := time.Now().Add(5 * time.Minute)
	claimsAccess := &models.Claims{
		Roles:    roles,
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			Subject:   user.Username,
			ExpiresAt: expirationTimeAccess.Unix(),
		},
	}
	// Создание токена на 5 часов
	expirationTimeRefresh := time.Now().Add(5 * time.Hour)
	claimsRefresh := &models.Claims{
		Roles:    roles,
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			Subject:   user.Username,
			ExpiresAt: expirationTimeRefresh.Unix(),
		},
	}

	tokenAccess := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsAccess)
	tokenRefresh := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefresh)

	tokenAccessStr, errAccess := tokenAccess.SignedString(jwtKey)
	if errAccess != nil {
		log.Println(errAccess.Error())
		c.JSON(500, gin.H{"message": "Cannot generate access token"})
		c.Abort()
		return
	}
	tokenRefreshStr, errRefresh := tokenRefresh.SignedString(jwtKey)
	if errRefresh != nil {
		log.Println(errRefresh.Error())
		c.JSON(500, gin.H{"message": "Cannot generate refresh token"})
		c.Abort()
		return
	}

	c.SetCookie("tokenAccess", tokenAccessStr, int(expirationTimeAccess.Unix()), "/", "localhost", false, true)
	c.SetCookie("tokenRefresh", tokenRefreshStr, int(expirationTimeRefresh.Unix()), "/", "localhost", false, true)
	c.JSON(200, gin.H{"message": "User logged in"})
}

func SignOut(c *gin.Context) {
	c.SetCookie("tokenAccess", "", -1, "/", "localhost", false, true)
	c.SetCookie("tokenRefresh", "", -1, "/", "localhost", false, true)
	c.JSON(200, gin.H{"message": "User logged out"})
}

func GetInfoAboutAccount(c *gin.Context) {
	var userInfo models.UserInfo

	cookie, err := c.Cookie("tokenAccess")
	if err != nil {
		c.JSON(400, gin.H{"message": "Unauthorized"})
		c.Abort()
		return
	}

	claims, err := helpers.ParseToken(cookie)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		c.Abort()
		return
	}

	row := databaseConn.QueryRow("SELECT * FROM Users WHERE username=$1", claims.Username)
	if err := row.Scan(&userInfo.UUID, &userInfo.Username, &userInfo.FirstName, &userInfo.LastName, &userInfo.Password); err != nil {
		c.JSON(400, gin.H{"message": "Cannot find user"})
		c.Abort()
		return
	}

	var roles []string
	rolesRows, err := databaseConn.Query("SELECT role FROM user_and_roles WHERE user_uuid=$1", userInfo.UUID)
	if err != nil {
		panic(err)
	}
	for rolesRows.Next() {
		var role string
		err := rolesRows.Scan(&role)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		roles = append(roles, role)
	}

	userInfo.Roles = roles

	c.JSON(200, userInfo)
}

func UpdateAccount(c *gin.Context) {
	var updateInfo models.UpdateUser

	if err := c.ShouldBindJSON(&updateInfo); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		c.Abort()
		return
	}

	cookie, err := c.Cookie("tokenAccess")
	if err != nil {
		c.JSON(400, gin.H{"message": "Unauthorized"})
		c.Abort()
		return
	}

	claims, err := helpers.ParseToken(cookie)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		c.Abort()
		return
	}

	_, err = databaseConn.Exec("UPDATE users SET last_name=$1, first_name=$2, password=$3 WHERE username=$4", updateInfo.LastName, updateInfo.FirstName, updateInfo.Password, claims.Username)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{"message": "User updated successfully"})
}