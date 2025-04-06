package storage

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/ruziba3vich/itv_test_project/internal/models"
	rediscl "github.com/ruziba3vich/itv_test_project/internal/redis_cl"
	"github.com/ruziba3vich/itv_test_project/internal/types"
	"gorm.io/gorm"
)

type MovieStorage struct {
	db            *gorm.DB
	redis_service *rediscl.RedisService
}

func NewMovieStorage(db *gorm.DB, redis_service *rediscl.RedisService) *MovieStorage {
	return &MovieStorage{db: db, redis_service: redis_service}
}

func (s *MovieStorage) Create(ctx context.Context, req *types.CreateMovieRequest) (*types.CreateMovieResponse, error) {
	movie := models.Movie{
		Title:    req.Title,
		Director: req.Director,
		Year:     req.Year,
		Plot:     req.Plot,
	}

	// Use a transaction for creating the movie
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&movie).Error; err != nil {
			return err
		}

		// Cache the movie in Redis within the transaction
		// If Redis fails, the whole operation fails
		if err := s.redis_service.SetMovie(ctx, &movie); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
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
	// Check Redis first
	movie, err := s.redis_service.GetMovie(ctx, req.ID)
	if err != nil && err != redis.Nil {
		// Only return error if it's not a cache miss
		return nil, err
	}

	if movie != nil {
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

	// Cache miss, fetch from DB
	movie = &models.Movie{}
	if err := s.db.WithContext(ctx).First(movie, req.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	// Cache in Redis before returning
	_ = s.redis_service.SetMovie(ctx, movie)

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

	// Use a transaction for updating the movie
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&movie, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("movie not found")
			}
			return err
		}

		// Update fields if provided
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

		if err := tx.Save(&movie).Error; err != nil {
			return err
		}

		// Update Redis cache within the transaction
		if err := s.redis_service.SetMovie(ctx, &movie); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
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

	// Use a transaction for deleting the movie
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&movie, req.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("movie not found")
			}
			return err
		}

		if err := tx.Delete(&movie).Error; err != nil {
			return err
		}

		// Remove from Redis within the transaction
		if err := s.redis_service.RemoveMovie(ctx, req.ID); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &types.DeleteMovieResponse{
		Message: "movie deleted successfully",
	}, nil
}
