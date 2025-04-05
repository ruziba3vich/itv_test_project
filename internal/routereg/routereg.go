package routereg

import (
	"github.com/gin-gonic/gin"
	handlers "github.com/ruziba3vich/itv_test_project/internal/http"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RegisterRoutes registers all routes, injecting the necessary dependencies
func RegisterRoutes(router *gin.Engine, authMiddleware func(gin.HandlerFunc) gin.HandlerFunc, handler *handlers.MovieHandler) {
	// Swagger route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Register your routes
	router.POST("/movies", authMiddleware(handler.CreateMovie))
	router.GET("/movies", handler.GetAllMovies)
	router.GET("/movies/:id", handler.GetMovieByID)
	router.PUT("/movies/:id", authMiddleware(handler.UpdateMovie))
	router.DELETE("/movies/:id", authMiddleware(handler.DeleteMovie))
}
