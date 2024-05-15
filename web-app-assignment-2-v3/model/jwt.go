package model

import "github.com/golang-jwt/jwt"

var JwtKey = []byte("secret-key")

type Claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}
