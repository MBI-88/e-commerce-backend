package middlewares

import (
	"fiber-backend/models"
	"fiber-backend/settings"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var keypass = settings.GetKeypass()

// JWTmiddleware verify token in headers
func JWTmiddleware(ctx *fiber.Ctx) error {
	tokenString := ctx.Get("Authorization", "")
	token, err := checkJWT(tokenString)
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Jwt error",
		})
	}
	if !token.Valid {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Invalid jwt",
		})
	}
	return ctx.Next()
}

// Generate JWT token
func JWTgenerate(user models.Users) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(15 * time.Minute).Unix(),
	})
	return token.SignedString(keypass)

}


// Admin jwt verify token in headers
func JWTadminMiddleware(ctx *fiber.Ctx) error {
	tokenString := ctx.Get("Authorization", "")
	token, err := checkJWT(tokenString)
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Jwt error",
		})
	}
	if !token.Valid {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Invalid jwt",
		})
	}
	
	if token.Claims.(*models.Claims).Role != "admin" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Bad credentials"})
	}
	return ctx.Next()
}


// JWT check helper
func checkJWT(tokenString string) (*jwt.Token, error) {
	tokenArray := strings.Fields(tokenString)
	if len(tokenArray) != 2 {
		return nil, fmt.Errorf("Lenght header error")
	}
	tokenCleaned := strings.TrimSpace(tokenArray[1])
	claims := models.Claims{}
	token, err := jwt.ParseWithClaims(tokenCleaned, &claims, func(j *jwt.Token) (interface{}, error) {
		return keypass, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}