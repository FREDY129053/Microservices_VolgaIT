package helpers

import (
	"document_microservice/database"
	"document_microservice/models"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func ParseToken(tokenStr string) (claims *models.Claims, err error) {
	token, err := jwt.ParseWithClaims(tokenStr, &models.Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(database.SECRET_KEY), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*models.Claims)
	if !ok {
		return nil, err
	}
	
	if time.Until(time.Unix(claims.ExpiresAt, 0)) < 0 {
		return nil, err
	}

	return claims, nil
}