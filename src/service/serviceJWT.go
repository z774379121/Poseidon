package service

import "github.com/dgrijalva/jwt-go"

type JwtCustomClaims struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}
