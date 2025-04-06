package types

import (
	"time"

	"github.com/ruziba3vich/itv_test_project/internal/models"
)

// CreateMovieRequest represents the request body for creating a movie
type (
	CreateMovieRequest struct {
		Title    string `json:"title" binding:"required"`
		Director string `json:"director" binding:"required"`
		Year     int    `json:"year" binding:"required,gte=1888,lte=2100"` // Reasonable year range
		Plot     string `json:"plot" binding:"max=1000"`                   // Optional, max length 1000 chars
	}

	// CreateMovieResponse represents the response after creating a movie
	CreateMovieResponse struct {
		ID        uint      `json:"id"`
		Title     string    `json:"title"`
		Director  string    `json:"director"`
		Year      int       `json:"year"`
		Plot      string    `json:"plot"`
		CreatedAt time.Time `json:"created_at"`
	}

	// GetAllRequest represents the query parameters for retrieving all movies
	GetAllRequest struct {
		Limit  int `json:"limit" form:"limit" binding:"min=1,max=100"` // Pagination limit
		Offset int `json:"offset" form:"offset" binding:"min=0"`       // Pagination offset
	}

	// GetAllResponse represents the response for retrieving all movies
	GetAllResponse struct {
		Movies     []models.Movie `json:"movies"`
		TotalCount int64          `json:"total_count"` // Total number of movies for pagination
	}

	// GetByIDRequest represents the request parameters for retrieving a movie by ID
	GetByIDRequest struct {
		ID uint `json:"id" uri:"id" binding:"required"`
	}

	// GetByIDResponse represents the response for retrieving a movie by ID
	GetByIDResponse struct {
		ID        uint      `json:"id"`
		Title     string    `json:"title"`
		Director  string    `json:"director"`
		Year      int       `json:"year"`
		Plot      string    `json:"plot"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	// UpdateMovieRequest represents the request body for updating a movie
	UpdateMovieRequest struct {
		Title    *string `json:"title" binding:"omitempty,min=1"`            // Optional, min length 1
		Director *string `json:"director" binding:"omitempty,min=1"`         // Optional
		Year     *int    `json:"year" binding:"omitempty,gte=1888,lte=2100"` // Optional
		Plot     *string `json:"plot" binding:"omitempty,max=1000"`          // Optional
	}

	// UpdateMovieResponse represents the response after updating a movie
	UpdateMovieResponse struct {
		ID        uint      `json:"id"`
		Title     string    `json:"title"`
		Director  string    `json:"director"`
		Year      int       `json:"year"`
		Plot      string    `json:"plot"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	// DeleteMovieRequest represents the request parameters for deleting a movie
	DeleteMovieRequest struct {
		ID uint `json:"id" uri:"id" binding:"required"`
	}

	// DeleteMovieResponse represents the response after deleting a movie
	DeleteMovieResponse struct {
		Message string `json:"message"`
	}

	CreateUserRequest struct {
		Fullname string `json:"full_name" binding:"required"`
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	LoginUserRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	LoginUserResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	RefreshTokenReq struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	RefreshTokenResponse struct {
		AccessToken string `json:"access_token"`
	}

	UsernameAlreadyTakenError struct {
		Message string `json:"message"`
	}
)

func (u *UsernameAlreadyTakenError) Error() string {
	return "this username is already taken"
}
