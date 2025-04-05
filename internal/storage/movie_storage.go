package storage

import (
	"context"
	"errors"
	"time"

	"github.com/ruziba3vich/itv_test_project/internal/models"
	"github.com/ruziba3vich/itv_test_project/internal/types"
	"gorm.io/gorm"
)

type MovieStorage struct {
	db *gorm.DB
}

func NewMovieStorage(db *gorm.DB) *MovieStorage {
	return &MovieStorage{db: db}
}

func (s *MovieStorage) Create(ctx context.Context, req *types.CreateMovieRequest) (*types.CreateMovieResponse, error) {
	movie := models.Movie{
		Title:    req.Title,
		Director: req.Director,
		Year:     req.Year,
		Plot:     req.Plot,
	}

	if err := s.db.WithContext(ctx).Create(&movie).Error; err != nil {
		return nil, err
	}

	return &types.CreateMovieResponse{
		ID:        movie.ID,
		Title:     movie.Title,
		Director:  movie.Director,
		Year:      movie.Year,
		Plot:      movie.Plot,
		CreatedAt: movie.CreatedAt,
	}, nil
}

func (s *MovieStorage) GetAll(ctx context.Context, req *types.GetAllRequest) (*types.GetAllResponse, error) {
	var (
		movies []models.Movie
		count  int64
	)

	if err := s.db.WithContext(ctx).Model(&models.Movie{}).Count(&count).Error; err != nil {
		return nil, err
	}

	if err := s.db.WithContext(ctx).
		Limit(req.Limit).
		Offset(req.Offset).
		Find(&movies).Error; err != nil {
		return nil, err
	}

	return &types.GetAllResponse{
		Movies:     movies,
		TotalCount: count,
	}, nil
}

func (s *MovieStorage) GetByID(ctx context.Context, req *types.GetByIDRequest) (*types.GetByIDResponse, error) {
	var movie models.Movie
	if err := s.db.WithContext(ctx).First(&movie, req.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &types.GetByIDResponse{
		ID:        movie.ID,
		Title:     movie.Title,
		Director:  movie.Director,
		Year:      movie.Year,
		Plot:      movie.Plot,
		CreatedAt: movie.CreatedAt,
		UpdatedAt: movie.UpdatedAt,
	}, nil
}

func (s *MovieStorage) Update(ctx context.Context, id uint, req *types.UpdateMovieRequest) (*types.UpdateMovieResponse, error) {
	var movie models.Movie

	if err := s.db.WithContext(ctx).First(&movie, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("movie not found")
		}
		return nil, err
	}

	// Only update fields that are not nil
	if req.Title != nil {
		movie.Title = *req.Title
	}
	if req.Director != nil {
		movie.Director = *req.Director
	}
	if req.Year != nil {
		movie.Year = *req.Year
	}
	if req.Plot != nil {
		movie.Plot = *req.Plot
	}
	movie.UpdatedAt = time.Now()

	if err := s.db.WithContext(ctx).Save(&movie).Error; err != nil {
		return nil, err
	}

	return &types.UpdateMovieResponse{
		ID:        movie.ID,
		Title:     movie.Title,
		Director:  movie.Director,
		Year:      movie.Year,
		Plot:      movie.Plot,
		UpdatedAt: movie.UpdatedAt,
	}, nil
}

func (s *MovieStorage) Delete(ctx context.Context, req *types.DeleteMovieRequest) (*types.DeleteMovieResponse, error) {
	var movie models.Movie
	if err := s.db.WithContext(ctx).First(&movie, req.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("movie not found")
		}
		return nil, err
	}

	if err := s.db.WithContext(ctx).Delete(&movie).Error; err != nil {
		return nil, err
	}

	return &types.DeleteMovieResponse{
		Message: "movie deleted successfully",
	}, nil
}
