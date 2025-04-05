package repos

import "context"

type (
	AuthRepo interface {
		GenerateTokens(ctx context.Context, userID uint) (string, string, error)
		RefreshAccessToken(ctx context.Context, refreshToken string) (string, error)
		ValidateJWT(tokenString string) (string, error)
	}
)
