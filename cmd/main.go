package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
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
			rl.NewTokenBucketLimiter,
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
