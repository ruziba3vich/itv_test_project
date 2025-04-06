package repos

import (
	"context"

	"github.com/ruziba3vich/itv_test_project/internal/models"
	"github.com/ruziba3vich/itv_test_project/internal/types"
)

type (
	AuthRepo interface {
		GenerateTokens(ctx context.Context, userID uint) (string, string, error)
		RefreshAccessToken(ctx context.Context, refreshToken string) (string, error)
		ValidateJWT(tokenString string) (string, error)
		LoginUser(ctx context.Context, req *types.LoginUserRequest) (uint, error)
		RegisterUser(ctx context.Context, user *models.User) error
	}
)
