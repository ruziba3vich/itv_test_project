package repos

import (
	"context"

	"github.com/ruziba3vich/itv_test_project/internal/models"
)

type IMovieService interface {
	Create(ctx context.Context, movie *models.Movie) error
	GetAll(ctx context.Context) ([]models.Movie, error)
	GetByID(ctx context.Context, id uint) (*models.Movie, error)
	Update(ctx context.Context, id uint, movie *models.Movie) error
	Delete(ctx context.Context, id uint) error
}
