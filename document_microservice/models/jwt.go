package models

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	Roles []string `json:"roles"`
	Username string `json:"username"`
	jwt.StandardClaims
}