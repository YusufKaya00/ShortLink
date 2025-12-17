package middleware

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ValidateResponse struct {
	Valid  bool      `json:"valid"`
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "bearer token required"})
			c.Abort()
			return
		}

		// Validate token with User Service
		userServiceURL := os.Getenv("USER_SERVICE_URL")
		if userServiceURL == "" {
			userServiceURL = "http://localhost:8081"
		}

		req, _ := http.NewRequest("GET", userServiceURL+"/api/users/validate", nil)
		req.Header.Set("Authorization", authHeader)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to validate token"})
			c.Abort()
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		var validateResp ValidateResponse
		if err := json.NewDecoder(resp.Body).Decode(&validateResp); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to parse validation response"})
			c.Abort()
			return
		}

		if !validateResp.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Set("user_id", validateResp.UserID)
		c.Set("email", validateResp.Email)
		c.Next()
	}
}

// OptionalAuthMiddleware tries to authenticate but doesn't require it
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.Next()
			return
		}

		userServiceURL := os.Getenv("USER_SERVICE_URL")
		if userServiceURL == "" {
			userServiceURL = "http://localhost:8081"
		}

		req, _ := http.NewRequest("GET", userServiceURL+"/api/users/validate", nil)
		req.Header.Set("Authorization", authHeader)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.Next()
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			c.Next()
			return
		}

		var validateResp ValidateResponse
		if err := json.NewDecoder(resp.Body).Decode(&validateResp); err != nil {
			c.Next()
			return
		}

		if validateResp.Valid {
			c.Set("user_id", validateResp.UserID)
			c.Set("email", validateResp.Email)
		}

		c.Next()
	}
}
