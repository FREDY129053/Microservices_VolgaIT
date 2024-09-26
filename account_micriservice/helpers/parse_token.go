package helpers

import (
	"account_microservice/models"
	"os"

	"github.com/dgrijalva/jwt-go"
)

func ParseToken(tokenStr string) (claims *models.Claims, err error) {
	token, err := jwt.ParseWithClaims(tokenStr, &models.Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*models.Claims)
	if !ok {
		return nil, err
	}

	return claims, nil
}