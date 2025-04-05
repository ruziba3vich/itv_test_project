package storage

import (
	"context"
	"errors"

	"github.com/ruziba3vich/itv_test_project/internal/models"
	"gorm.io/gorm"
)

// MovieStorage is the concrete implementation of the Storage interface
type MovieStorage struct {
	db *gorm.DB
}

// NewMovieStorage creates a new instance of MovieStorage
func NewMovieStorage(db *gorm.DB) *MovieStorage {
	return &MovieStorage{db: db}
}

// Create adds a new movie to the database
func (s *MovieStorage) Create(ctx context.Context, movie *models.Movie) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(movie).Error; err != nil {
			return err
		}
		return nil
	})
}

// GetAll retrieves all movies from the database
func (s *MovieStorage) GetAll(ctx context.Context, limit, offset int) ([]models.Movie, int64, error) {
	var (
		movies []models.Movie
		count  int64
	)

	// Count total
	if err := s.db.WithContext(ctx).Model(&models.Movie{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Paginate
	if err := s.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Find(&movies).Error; err != nil {
		return nil, 0, err
	}

	return movies, count, nil
}

// GetByID retrieves a movie by its ID
func (s *MovieStorage) GetByID(ctx context.Context, id uint) (*models.Movie, error) {
	var movie models.Movie
	if err := s.db.WithContext(ctx).First(&movie, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if not found, no error
		}
		return nil, err
	}
	return &movie, nil
}

// Update modifies an existing movie by ID
func (s *MovieStorage) Update(ctx context.Context, id uint, movie *models.Movie) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existingMovie models.Movie
		if err := tx.First(&existingMovie, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("movie not found")
			}
			return err
		}

		movie.ID = id
		if err := tx.Model(&existingMovie).Updates(movie).Error; err != nil {
			return err
		}
		return nil
	})
}

// Delete removes a movie by ID (soft delete via GORM)
func (s *MovieStorage) Delete(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var movie models.Movie
		if err := tx.First(&movie, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("movie not found")
			}
			return err
		}

		if err := tx.Delete(&movie).Error; err != nil {
			return err
		}
		return nil
	})
}
