package service

import (
	"context"

	"github.com/ruziba3vich/itv_test_project/internal/storage"
	"github.com/ruziba3vich/itv_test_project/internal/types"
	"github.com/ruziba3vich/itv_test_project/pkg/logger"
)

// MovieService represents the service layer for movies
type MovieService struct {
	storage *storage.MovieStorage
	logger  *logger.Logger
}

// NewMovieService initializes a new MovieService
func NewMovieService(storage *storage.MovieStorage, logger *logger.Logger) *MovieService {
	return &MovieService{storage: storage, logger: logger}
}

// CreateMovie creates a new movie
func (s *MovieService) CreateMovie(ctx context.Context, req *types.CreateMovieRequest) (*types.CreateMovieResponse, error) {
	// Call the storage layer to create the movie
	resp, err := s.storage.Create(ctx, req)
	if err != nil {
		// Log the error in the service layer
		s.logger.Error("Failed to create movie", map[string]any{
			"title":    req.Title,
			"director": req.Director,
			"year":     req.Year,
			"error":    err.Error(),
		})
		return nil, err
	}
	return resp, nil
}

// DeleteMovie deletes a movie by ID
func (s *MovieService) DeleteMovie(ctx context.Context, req *types.DeleteMovieRequest) (*types.DeleteMovieResponse, error) {
	// Call the storage layer to delete the movie
	resp, err := s.storage.Delete(ctx, req)
	if err != nil {
		// Log the error
		s.logger.Error("Failed to delete movie", map[string]any{
			"id":    req.ID,
			"error": err.Error(),
		})
		return nil, err
	}
	return resp, nil
}

// GetAllMovies retrieves all movies with pagination
func (s *MovieService) GetAllMovies(ctx context.Context, req *types.GetAllRequest) (*types.GetAllResponse, error) {
	// Call the storage layer to get the movies
	resp, err := s.storage.GetAll(ctx, req)
	if err != nil {
		// Log the error
		s.logger.Error("Failed to retrieve all movies", map[string]any{
			"limit":  req.Limit,
			"offset": req.Offset,
			"error":  err.Error(),
		})
		return nil, err
	}
	return resp, nil
}

// GetMovieByID retrieves a movie by its ID
func (s *MovieService) GetMovieByID(ctx context.Context, req *types.GetByIDRequest) (*types.GetByIDResponse, error) {
	// Call the storage layer to get the movie by ID
	resp, err := s.storage.GetByID(ctx, req)
	if err != nil {
		// Log the error
		s.logger.Error("Failed to retrieve movie by ID", map[string]any{
			"id":    req.ID,
			"error": err.Error(),
		})
		return nil, err
	}
	return resp, nil
}

// UpdateMovie updates an existing movie by ID
func (s *MovieService) UpdateMovie(ctx context.Context, id uint, req *types.UpdateMovieRequest) (*types.UpdateMovieResponse, error) {
	// Call the storage layer to update the movie
	resp, err := s.storage.Update(ctx, id, req)
	if err != nil {
		// Log the error
		s.logger.Error("Failed to update movie", map[string]any{
			"id":    id,
			"error": err.Error(),
		})
		return nil, err
	}
	return resp, nil
}
