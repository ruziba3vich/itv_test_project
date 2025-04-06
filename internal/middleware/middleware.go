package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ruziba3vich/itv_test_project/internal/repos"
	"github.com/ruziba3vich/itv_test_project/pkg/logger"
	limiter "github.com/ruziba3vich/prodonik_rl"
)

// AuthHandler holds dependencies for authentication
type AuthHandler struct {
	authRepo repos.AuthRepo
	logger   *logger.Logger
	limiter  *limiter.TokenBucketLimiter
}

// NewAuthHandler initializes and returns an AuthHandler instance
func NewAuthHandler(authRepo repos.AuthRepo, logger *logger.Logger, limiter *limiter.TokenBucketLimiter) *AuthHandler {
	return &AuthHandler{
		authRepo: authRepo,
		logger:   logger,
		limiter:  limiter,
	}
}

// AuthMiddleware validates JWT and sets user ID before executing the given handlers
func (a *AuthHandler) AuthMiddleware() func(gin.HandlerFunc) gin.HandlerFunc {
	return func(handler gin.HandlerFunc) gin.HandlerFunc {
		return func(c *gin.Context) {
			ip := c.ClientIP() // Get user IP for rate limiting

			allowed, err := a.limiter.AllowRequest(c, ip)
			if err != nil {
				a.logger.Println("Rate limiter error:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				c.Abort()
				return
			}

			if !allowed {
				a.logger.Println("Rate limit exceeded for IP:", ip)
				c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
				c.Abort()
				return
			}

			tokenString := c.GetHeader("Authorization")
			if tokenString == "" {
				a.logger.Println("Missing authorization token")
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
				c.Abort()
				return
			}

			parts := strings.Split(tokenString, " ")

			userID, err := a.authRepo.ValidateJWT(parts[1])
			if err != nil {
				a.logger.Println("Invalid token:", err)
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
				c.Abort()
				return
			}

			// Set user ID in context
			c.Set("userID", userID)

			// Call the actual handler
			handler(c)
		}
	}
}
