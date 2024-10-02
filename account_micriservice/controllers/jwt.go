package controllers

import (
	"github.com/gin-gonic/gin"
	"account_microservice/helpers"
	"account_microservice/models"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)


// VerifyingToken godoc
// VerifyingToken верификация access токена
// @Summary Проверка access токена
// @Description Проверка access токена на поддельность
// @Tags JWT
// @Accept json
// @Produce json
// @Param access_token path string true "Access токен"
// @Success 200 {object} []string "Token verified successfully"
// @Failure 400 {object} map[string]string "Token verification failed"
// @Router /Authentication/Validate [get]
// @Security ApiKeyAuth
func VerifyingToken(c *gin.Context) {
	access_token := c.Query("access_token")

	_, err := helpers.VerifyToken(access_token)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{"message": "Token verification failed"})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{"message": "Token verified successfully"})
}


// RefreshAccessToken godoc
// RefreshAccessToken обновление access токена
// @Summary Обновление access токена
// @Description Обновление access токена с помощью refresh токена
// @Tags JWT
// @Accept json
// @Produce json
// @Success 200 {object} []string "Token refreshed successfully"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /Authentication/Refresh [get]
// @Security ApiKeyAuth
func RefreshAccessToken(c *gin.Context) {
	refreshToken, err := c.Cookie("tokenRefresh")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	

	claims, err := helpers.ParseToken(refreshToken)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		c.Abort()
		return
	}

	expirationTimeAccess := time.Now().Add(5 * time.Minute)
	claimsNewAccess := &models.Claims{
		Roles:    claims.Roles,
		Username: claims.Username,
		StandardClaims: jwt.StandardClaims{
			Subject:   claims.Username,
			ExpiresAt: expirationTimeAccess.Unix(),
		},
	}
	tokenNewAccess := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsNewAccess)
	tokenNewAccessStr, errAccess := tokenNewAccess.SignedString(jwtKey)
	if errAccess != nil {
		log.Println(errAccess.Error())
		c.JSON(500, gin.H{"message": "Cannot refresh access token"})
		c.Abort()
		return
	}
	c.SetCookie("tokenAccess", tokenNewAccessStr, int(expirationTimeAccess.Unix()), "/", "localhost", false, true)
	

	c.JSON(200, gin.H{"message": "Token refreshed successfully"})
}