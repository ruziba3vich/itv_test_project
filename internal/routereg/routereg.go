package routereg

import (
	"github.com/gin-gonic/gin"
	handlers "github.com/ruziba3vich/itv_test_project/internal/http"
	"github.com/ruziba3vich/itv_test_project/internal/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RegisterRoutes registers all routes, injecting the necessary dependencies
func RegisterMovieRoutes(router *gin.Engine, middleware *middleware.AuthHandler, handler *handlers.MovieHandler) {
	// Swagger route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	authMiddleware := middleware.AuthMiddleware()
	// Register your routes
	router.POST("/movies", authMiddleware(handler.CreateMovie))
	router.GET("/movies", handler.GetAllMovies)
	router.GET("/movies/:id", handler.GetMovieByID)
	router.PUT("/movies/:id", authMiddleware(handler.UpdateMovie))
	router.DELETE("/movies/:id", authMiddleware(handler.DeleteMovie))
}

// RegisterRoutes registers all authentication-related routes
func RegisterAuthRoutes(router *gin.Engine, handler *handlers.AuthHandler) {
	router.POST("/register", handler.RegisterUser) // Separate endpoint for registration
	router.POST("/login", handler.Login)
	router.POST("/refresh", handler.RefreshToken)
}
