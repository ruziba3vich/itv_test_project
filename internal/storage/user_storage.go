package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/ruziba3vich/itv_test_project/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type (
	UserStorage struct {
		db *gorm.DB
	}
)

// CreateUser adds a new user to the database
func (s *UserStorage) CreateUser(ctx context.Context, user *models.User) error {
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	existingUser, err := s.getUserByUsername(ctx, *user.Username)
	if err != nil {
		return fmt.Errorf("failed to verify username availability: %s", err.Error())
	}
	if existingUser != nil {
		return fmt.Errorf("this username is already taken")
	}
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		return nil
	})
}

// GetUserByUsername retrieves a user by username
func (s *UserStorage) getUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	if err := s.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // User not found
		}
		return nil, err
	}
	return &user, nil
}

// CreateRefreshToken stores a new refresh token
func (s *UserStorage) CreateRefreshToken(ctx context.Context, token *models.RefreshToken) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(token).Error; err != nil {
			return err
		}
		return nil
	})
}

// GetRefreshToken retrieves a refresh token by its value
func (s *UserStorage) GetRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	if err := s.db.Where("token = ?", token).First(&refreshToken).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Token not found
		}
		return nil, err
	}
	return &refreshToken, nil
}

// Login checks user credentials and returns a JWT token
func (s *UserStorage) Login(ctx context.Context, username, password string) error {
	user, err := s.getUserByUsername(ctx, username)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Check password
	if !checkPassword(user.Password, password) {
		return errors.New("invalid username or password")
	}

	return nil
}

// hashPassword hashes a password using bcrypt
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword verifies a hashed password
func checkPassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
