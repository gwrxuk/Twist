package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/twist/api-gateway/internal/models"
)

// JWTClaims represents the claims in the JWT
type JWTClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// Auth middleware validates the JWT token
func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.NewErrorResponse("Authorization header is required"))
			return
		}

		// Check if it's a Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.NewErrorResponse("Invalid authorization format, Bearer token required"))
			return
		}

		// Extract the token
		tokenString := parts[1]

		// Parse the token
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.NewErrorResponse("Invalid or expired token"))
			return
		}

		// Extract claims
		claims, ok := token.Claims.(*JWTClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.NewErrorResponse("Invalid token claims"))
			return
		}

		// Check if token is expired
		if time.Until(claims.ExpiresAt.Time) < 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.NewErrorResponse("Token expired"))
			return
		}

		// Set user data in the context
		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.NewErrorResponse("Invalid user ID in token"))
			return
		}

		c.Set("user_id", userID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}
