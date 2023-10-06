package models

import "github.com/golang-jwt/jwt/v5"

// JWTClient. Client for getting user's data
type JWTClient struct {
	Name  string `json:"username"`
	Token string `json:"token"`
}

// Claims. User's data to put into token
type Claims struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}
