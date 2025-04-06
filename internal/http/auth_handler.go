// @title ITV Test Project API
// @version 1.0
// @description API for movie management
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:7777
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the JWT token.

package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ruziba3vich/itv_test_project/internal/models"
	"github.com/ruziba3vich/itv_test_project/internal/repos"
	"github.com/ruziba3vich/itv_test_project/internal/types"
	"github.com/ruziba3vich/itv_test_project/pkg/logger"
)

// AuthHandler manages authentication-related endpoints
type AuthHandler struct {
	authRepo repos.AuthRepo // Abstract field for token operations
	log      *logger.Logger // For logging
}

// NewAuthHandler creates a new AuthHandler with dependencies
func NewAuthHandler(authRepo repos.AuthRepo, log *logger.Logger) *AuthHandler {
	return &AuthHandler{
		authRepo: authRepo,
		log:      log,
	}
}

// RegisterUser godoc
// @Summary Register a new user
// @Description Registers a new user with the provided credentials
// @Tags auth
// @Accept json
// @Produce json
// @Param user body types.CreateUserRequest true "User registration data"
// @Success 200 {object} gin.H "message: User registered successfully"
// @Failure 400 {object} gin.H "error: invalid request"
// @Failure 409 {object} gin.H "error: username already taken"
// @Failure 500 {object} gin.H "error: registration failed"
// @Router /register [post]
func (h *AuthHandler) RegisterUser(c *gin.Context) {
	var req types.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("Invalid registration request", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := h.authRepo.RegisterUser(c.Request.Context(), &models.User{
		Fullname: req.Fullname,
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, &types.UsernameAlreadyTakenError{}) {
			h.log.Warn("Duplicate username during registration", map[string]any{
				"username": req.Username,
			})
			c.JSON(http.StatusConflict, gin.H{"error": "username already taken"})
			return
		}
		h.log.Error("Failed to register user", map[string]any{
			"error":    err.Error(),
			"username": req.Username,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "registration failed: " + err.Error()})
		return
	}

	h.log.Info("User registered successfully", map[string]any{
		"username": req.Username,
	})
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// Login godoc
// @Summary User login
// @Description Authenticates a user and returns access and refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body types.LoginUserRequest true "User login credentials"
// @Success 200 {object} types.LoginUserResponse "Access and refresh tokens"
// @Failure 400 {object} gin.H "error: invalid request"
// @Failure 401 {object} gin.H "error: invalid credentials"
// @Failure 500 {object} gin.H "error: login failed"
// @Router /login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req types.LoginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("Invalid login request", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	id, err := h.authRepo.LoginUser(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "invalid credentials" {
			h.log.Warn("Invalid login attempt", map[string]interface{}{
				"username": req.Username,
			})
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		h.log.Error("Failed to login user", map[string]interface{}{
			"error":    err.Error(),
			"username": req.Username,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "login failed"})
		return
	}

	accessTokenStr, refreshTokenStr, err := h.authRepo.GenerateTokens(c.Request.Context(), id)
	if err != nil {
		h.log.Error("Failed to generate tokens", map[string]interface{}{
			"error":   err.Error(),
			"user_id": id,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
		return
	}

	h.log.Info("User logged in successfully", map[string]interface{}{
		"user_id": id,
	})
	c.JSON(http.StatusOK, types.LoginUserResponse{
		AccessToken:  accessTokenStr,
		RefreshToken: refreshTokenStr,
	})
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Generates a new access token using a refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body types.RefreshTokenReq true "Refresh token request"
// @Success 200 {object} types.RefreshTokenResponse "New access token"
// @Failure 400 {object} gin.H "error: invalid request"
// @Failure 401 {object} gin.H "error: invalid or expired refresh token"
// @Failure 500 {object} gin.H "error: refresh failed"
// @Router /refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req types.RefreshTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("Invalid refresh token request", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	accessToken, err := h.authRepo.RefreshAccessToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		h.log.Warn("Failed to refresh access token", map[string]interface{}{
			"error": err.Error(),
			"token": req.RefreshToken,
		})
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired refresh token"})
		return
	}

	h.log.Info("Access token refreshed successfully", map[string]interface{}{
		"token": req.RefreshToken,
	})
	c.JSON(http.StatusOK, types.RefreshTokenResponse{
		AccessToken: accessToken,
	})
}
