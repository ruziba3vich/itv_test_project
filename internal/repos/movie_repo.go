package repos

import (
	"context"

	"github.com/ruziba3vich/itv_test_project/internal/types"
)

type IMovieService interface {
	Create(ctx context.Context, req *types.CreateMovieRequest) (*types.CreateMovieResponse, error)
	Delete(ctx context.Context, req *types.DeleteMovieRequest) (*types.DeleteMovieResponse, error)
	GetAll(ctx context.Context, req *types.GetAllRequest) (*types.GetAllResponse, error)
	GetByID(ctx context.Context, req *types.GetByIDRequest) (*types.GetByIDResponse, error)
	Update(ctx context.Context, id uint, req *types.UpdateMovieRequest) (*types.UpdateMovieResponse, error)
}
