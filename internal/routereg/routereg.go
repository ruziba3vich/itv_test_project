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

// @title ITV Test Project API
// @version 1.0
// @description API for movie management
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:7777
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the JWT token.
func SwaggerInfo() {}

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
