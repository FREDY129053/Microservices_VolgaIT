// Claims и модели для JWT
package models

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	Roles []string `json:"roles"`
	jwt.StandardClaims
}

type ValidateToken struct {
	AccessToken string `json:"access_token"`
}

type RefreshToken struct {
	RefreshToken string `json:"refresh_token"`
}