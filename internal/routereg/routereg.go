package routereg

import (
	"github.com/gin-gonic/gin"
	_ "github.com/ruziba3vich/itv_test_project/docs"
	handlers "github.com/ruziba3vich/itv_test_project/internal/http"
	"github.com/ruziba3vich/itv_test_project/internal/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/swaggo/swag"
)

// RegisterRoutes registers all routes, injecting the necessary dependencies
func RegisterMovieRoutes(router *gin.Engine, middleware *middleware.AuthHandler, handler *handlers.MovieHandler) {
	// Swagger route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	authMiddleware := middleware.AuthMiddleware()
	movie_router := router.Group("api/v1")
	// Register your routes
	movie_router.POST("/movies", authMiddleware(handler.CreateMovie))
	movie_router.GET("/movies", handler.GetAllMovies)
	movie_router.GET("/movies/:id", handler.GetMovieByID)
	movie_router.PUT("/movies/:id", authMiddleware(handler.UpdateMovie))
	movie_router.DELETE("/movies/:id", authMiddleware(handler.DeleteMovie))
}

// RegisterRoutes registers all authentication-related routes
func RegisterAuthRoutes(router *gin.Engine, handler *handlers.AuthHandler) {
	movie_router := router.Group("api/v1")
	movie_router.POST("/register", handler.RegisterUser) // Separate endpoint for registration
	movie_router.POST("/login", handler.Login)
	movie_router.POST("/refresh", handler.RefreshToken)
}
