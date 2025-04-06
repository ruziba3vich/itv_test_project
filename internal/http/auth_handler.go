package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ruziba3vich/itv_test_project/internal/repos"
	"github.com/ruziba3vich/itv_test_project/internal/service"
	"github.com/ruziba3vich/itv_test_project/internal/types"
	"github.com/ruziba3vich/itv_test_project/pkg/logger"
)

// AuthHandler manages authentication-related endpoints
type AuthHandler struct {
	authRepo repos.AuthRepo        // Abstract field for token operations
	userSvc  *service.TokenService // For user validation during login
	log      *logger.Logger        // For logging
}

// NewAuthHandler creates a new AuthHandler with dependencies
func NewAuthHandler(authRepo repos.AuthRepo, userSvc *service.TokenService, log *logger.Logger) *AuthHandler {
	return &AuthHandler{
		authRepo: authRepo,
		userSvc:  userSvc,
		log:      log,
	}
}

// RegisterRoutes sets up authentication routes
func (h *AuthHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/login", h.Login)
	router.POST("/refresh", h.RefreshToken)
}

// Login godoc
// @Summary User login
// @Description Authenticates a user and returns access and refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body struct{ Username string; Password string } true "User credentials"
// @Success 200 {object} map[string]string
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 500 {object} gin.H
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
		h.log.Error("Failed to retrieve user", map[string]interface{}{
			"error":    err.Error(),
			"username": req.Username,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "login failed"})
		return
	}

	accessTokenStr, refreshTokenStr, err := h.authRepo.GenerateTokens(c.Request.Context(), id)
	if err != nil {
		h.log.Error("Failed to retrieve user", map[string]interface{}{
			"error":    err.Error(),
			"username": req.Username,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "login failed"})
		return
	}

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
// @Param refresh_token body struct{ RefreshToken string } true "Refresh token"
// @Success 200 {object} map[string]string
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 500 {object} gin.H
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

	// Refresh the access token
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
