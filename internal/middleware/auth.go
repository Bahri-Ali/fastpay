package middleware

import (
    "fastpay-backend/internal/auth"
    "fastpay-backend/pkg/utils"
    "net/http"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
)

func AuthMiddleware(authRepo auth.Repository) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }

        parts := strings.Split(authHeader, "Bearer ")
        if len(parts) != 2 {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
            c.Abort()
            return
        }
        rawToken := parts[1]

        tokenHash := utils.HashToken(rawToken)

        session, err := authRepo.GetSessionByTokenHash(c.Request.Context(), tokenHash)
        if err != nil || session == nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
            c.Abort()
            return
        }

        if time.Now().After(session.ExpiresAt) {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Session expired"})
            c.Abort()
            return
        }

        newExpiry := time.Now().Add(48 * time.Hour)
        if err := authRepo.UpdateSessionExpiry(c.Request.Context(), tokenHash, newExpiry); err != nil {
        }

        c.Set("user_id", session.UserID)
        c.Next()
    }
}