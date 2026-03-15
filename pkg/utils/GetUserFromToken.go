package utils

import (
	"strings"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gin-gonic/gin"
)

var jwtSecret = []byte("your_secret_key")

func GetUserFromToken(c *gin.Context) (bool, int64) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		return false, 0
	}

	// Expect: "Bearer <token>"
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return false, 0
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, 0
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return false, 0
	}

	return true, int64(userIDFloat)
}