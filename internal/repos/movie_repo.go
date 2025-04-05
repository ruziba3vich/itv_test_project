package repos

import (
	"context"

	"github.com/ruziba3vich/itv_test_project/internal/types"
)

type IMovieService interface {
	CreateMovie(ctx context.Context, req *types.CreateMovieRequest) (*types.CreateMovieResponse, error)
	DeleteMovie(ctx context.Context, req *types.DeleteMovieRequest) (*types.DeleteMovieResponse, error)
	GetAllMovies(ctx context.Context, req *types.GetAllRequest) (*types.GetAllResponse, error)
	GetMovieByID(ctx context.Context, req *types.GetByIDRequest) (*types.GetByIDResponse, error)
	UpdateMovie(ctx context.Context, id uint, req *types.UpdateMovieRequest) (*types.UpdateMovieResponse, error)
}
