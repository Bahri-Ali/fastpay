package middleware

import (
	"fastpay-backend/internal/auth"
	"fastpay-backend/pkg/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// we are dont use this middleware for new last update 23/3 by bahri ali 
func WSAuthMiddleware(authRepo auth.Repository) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(401, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }

        parts := strings.Split(authHeader, "Bearer ")
        if len(parts) != 2 {
            c.JSON(401, gin.H{"error": "Invalid authorization format"})
            c.Abort()
            return
        }

        rawToken := parts[1]
        tokenHash := utils.HashToken(rawToken)

        session, err := authRepo.GetSessionByTokenHash(c.Request.Context(), tokenHash)
        if err != nil || session == nil {
            c.JSON(401, gin.H{"error": "Invalid session"})
            c.Abort()
            return
        }

        if time.Now().After(session.ExpiresAt) {
            c.JSON(401, gin.H{"error": "Session expired"})
            c.Abort()
            return
        }

        c.Set("user_id", session.UserID)

        c.Next()
    }
}