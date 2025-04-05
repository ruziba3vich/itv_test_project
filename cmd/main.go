package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	handlers "github.com/ruziba3vich/itv_test_project/internal/http"
	"github.com/ruziba3vich/itv_test_project/internal/middleware"
	"github.com/ruziba3vich/itv_test_project/internal/routereg"
	"github.com/ruziba3vich/itv_test_project/pkg/config"
	"github.com/ruziba3vich/itv_test_project/pkg/db"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Provide(
			config.LoadDBConfig,
			db.NewDB,
			NewGinEngine,
			handlers.NewMovieHandler,
			middleware.NewAuthHandler,
			
		),
		fx.Invoke(
			routereg.RegisterRoutes,
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
