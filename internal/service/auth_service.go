package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ruziba3vich/itv_test_project/internal/models"
	"github.com/ruziba3vich/itv_test_project/internal/storage"
	"github.com/ruziba3vich/itv_test_project/pkg/logger"
)

// TokenService implementation
type TokenService struct {
	store      *storage.UserStorage
	log        *logger.Logger
	secret     string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

// NewTokenService creates a new TokenService
func NewTokenService(store *storage.UserStorage, log *logger.Logger, secret string) *TokenService {
	return &TokenService{
		store:      store,
		log:        log,
		secret:     secret,
		accessTTL:  15 * time.Minute,   // Access token expires in 15 minutes
		refreshTTL: 7 * 24 * time.Hour, // Refresh token expires in 7 days
	}
}

// GenerateTokens creates an access token and refresh token for a user
func (s *TokenService) GenerateTokens(ctx context.Context, userID uint) (string, string, error) {
	// Generate access token
	accessClaims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(s.accessTTL).Unix(),
		"iat": time.Now().Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenStr, err := accessToken.SignedString([]byte(s.secret))
	if err != nil {
		s.log.Error("Failed to generate access token", map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
		})
		return "", "", err
	}

	// Generate refresh token
	refreshTokenStr := uuid.New().String() // Random UUID as refresh token
	refreshToken := &models.RefreshToken{
		UserID:    userID,
		Token:     refreshTokenStr,
		ExpiresAt: time.Now().Add(s.refreshTTL),
	}
	if err := s.store.CreateRefreshToken(ctx, refreshToken); err != nil {
		return "", "", err
	}

	s.log.Info("Tokens generated successfully", map[string]interface{}{
		"user_id": userID,
	})
	return accessTokenStr, refreshTokenStr, nil
}

// RefreshAccessToken generates a new access token using a valid refresh token
func (s *TokenService) RefreshAccessToken(ctx context.Context, refreshToken string) (string, error) {
	// Validate refresh token
	rt, err := s.store.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", err
	}
	if rt == nil || rt.ExpiresAt.Before(time.Now()) {
		s.log.Warn("Invalid or expired refresh token", map[string]interface{}{
			"token": refreshToken,
		})
		return "", errors.New("invalid or expired refresh token")
	}

	// Generate new access token
	accessClaims := jwt.MapClaims{
		"sub": rt.UserID,
		"exp": time.Now().Add(s.accessTTL).Unix(),
		"iat": time.Now().Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenStr, err := accessToken.SignedString([]byte(s.secret))
	if err != nil {
		s.log.Error("Failed to generate new access token", map[string]interface{}{
			"error":   err.Error(),
			"user_id": rt.UserID,
		})
		return "", err
	}

	s.log.Info("Access token refreshed", map[string]interface{}{
		"user_id": rt.UserID,
	})
	return accessTokenStr, nil
}

// ValidateJWT validates a JWT token and returns the user ID
func (s *TokenService) ValidateJWT(tokenString string) (string, error) {
	// Parse the token and verify the signature method
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// Check that the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secret), nil
	})

	// Return an error if the token could not be parsed
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %v", err)
	}

	// Extract and validate claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, exists := claims["user_id"].(string); exists {
			return userID, nil
		}
		return "", fmt.Errorf("user_id not found in token claims")
	}

	// Return an error if the token is not valid
	return "", fmt.Errorf("invalid token")
}

// RegisterUser creates a new user
func (s *TokenService) RegisterUser(ctx context.Context, user *models.User) error {
	err := s.store.CreateUser(ctx, user)
	if err != nil {
		s.log.Error("Failed to create user: " + err.Error())
	}
	return err
}
