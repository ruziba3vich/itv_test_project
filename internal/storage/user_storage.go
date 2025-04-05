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

// hashPassword hashes a password using bcrypt
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
