package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	handlers "github.com/ruziba3vich/itv_test_project/internal/http"
	"github.com/ruziba3vich/itv_test_project/internal/middleware"
	"github.com/ruziba3vich/itv_test_project/internal/routereg"
	"github.com/ruziba3vich/itv_test_project/internal/service"
	"github.com/ruziba3vich/itv_test_project/internal/storage"
	"github.com/ruziba3vich/itv_test_project/pkg/config"
	"github.com/ruziba3vich/itv_test_project/pkg/db"
	"github.com/ruziba3vich/itv_test_project/pkg/logger"
	"github.com/ruziba3vich/itv_test_project/pkg/rediscl"
	rl "github.com/ruziba3vich/prodonik_rl"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Provide(
			config.LoadConfig,
			rediscl.NewRedisClient,
			logger.NewLogger,
			db.NewDB,
			NewRateLimiter,
			storage.NewMovieStorage,
			storage.NewUserStorage,
			service.NewMovieService,
			service.NewTokenService,
			NewGinEngine,
			handlers.NewMovieHandler,
			handlers.NewAuthHandler,
			middleware.NewAuthHandler,
		),
		fx.Invoke(
			routereg.RegisterMovieRoutes,
			routereg.RegisterAuthRoutes,
			RunServer, // Add this new function to start the server
		),
	)

	// Start the application
	if err := app.Start(context.Background()); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}

	// Wait for the app to stop (e.g., via SIGTERM)
	<-app.Done()
	log.Println("Application stopped")
}

// NewGinEngine provides the Gin engine instance
func NewGinEngine() *gin.Engine {
	return gin.Default() // This creates a new Gin engine instance with default middleware
}

// RunServer starts the Gin server on port 7777
func RunServer(lc fx.Lifecycle, router *gin.Engine, logger *logger.Logger, cfg *config.Config) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				logger.Info("Starting server on port 7777")
				if err := router.Run(":" + cfg.AppPort); err != nil {
					logger.Error("Failed to start server: " + err.Error())
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Shutting down server")
			return nil
		},
	})
}

func NewRateLimiter(redisClient *redis.Client, cfg *config.Config) *rl.TokenBucketLimiter {
	return rl.NewTokenBucketLimiter(redisClient, cfg.RLConfig.MaxTokens, float64(cfg.RLConfig.RefillRate), cfg.RLConfig.Window)
}
