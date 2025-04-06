package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ruziba3vich/itv_test_project/internal/repos"
	"github.com/ruziba3vich/itv_test_project/internal/service"
	"github.com/ruziba3vich/itv_test_project/internal/types"
	"github.com/ruziba3vich/itv_test_project/pkg/logger"

	_ "github.com/swaggo/files"       // Swagger UI files
	_ "github.com/swaggo/gin-swagger" // Required for Swagger
)

// MovieHandler handles HTTP requests for movies
type MovieHandler struct {
	svc repos.IMovieService
	log *logger.Logger
}

// NewMovieHandler creates a new MovieHandler with dependencies
func NewMovieHandler(svc *service.MovieService, log *logger.Logger) *MovieHandler {
	return &MovieHandler{svc: svc, log: log}
}

// CreateMovie godoc
// @Summary Create a new movie
// @Description Creates a new movie record in the database
// @Tags movies
// @Accept json
// @Produce json
// @Param movie body types.CreateMovieRequest true "Movie data"
// @Success 201 {object} types.CreateMovieResponse
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Security BearerAuth
// @Router /movies [post]
func (h *MovieHandler) CreateMovie(c *gin.Context) {
	var req types.CreateMovieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("Invalid create movie request", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	resp, err := h.svc.CreateMovie(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create movie"})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GetAllMovies godoc
// @Summary Get all movies
// @Description Retrieves a paginated list of all movies
// @Tags movies
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} types.GetAllResponse
// @Failure 500 {object} gin.H
// @Security BearerAuth
// @Router /movies [get]
func (h *MovieHandler) GetAllMovies(c *gin.Context) {
	var req types.GetAllRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.log.Warn("Invalid get all movies request", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Set defaults if not provided
	if req.Limit == 0 {
		req.Limit = 10
	}

	resp, err := h.svc.GetAllMovies(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve movies"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetMovieByID godoc
// @Summary Get a movie by ID
// @Description Retrieves a specific movie by its ID
// @Tags movies
// @Produce json
// @Param id path int true "Movie ID"
// @Success 200 {object} types.GetByIDResponse
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Security BearerAuth
// @Router /movies/{id} [get]
func (h *MovieHandler) GetMovieByID(c *gin.Context) {
	var req types.GetByIDRequest
	if err := c.ShouldBindUri(&req); err != nil {
		h.log.Warn("Invalid get movie by ID request", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	resp, err := h.svc.GetMovieByID(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get movie"})
		return
	}
	if resp == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateMovie godoc
// @Summary Update a movie
// @Description Updates an existing movie by ID
// @Tags movies
// @Accept json
// @Produce json
// @Param id path int true "Movie ID"
// @Param movie body types.UpdateMovieRequest true "Updated movie data"
// @Success 200 {object} types.UpdateMovieResponse
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Security BearerAuth
// @Router /movies/{id} [put]
func (h *MovieHandler) UpdateMovie(c *gin.Context) {
	var req types.UpdateMovieRequest
	var idReq types.GetByIDRequest

	if err := c.ShouldBindUri(&idReq); err != nil {
		h.log.Warn("Invalid update movie ID", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("Invalid update movie request", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	resp, err := h.svc.UpdateMovie(c.Request.Context(), idReq.ID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update movie"})
		return
	}
	if resp == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteMovie godoc
// @Summary Delete a movie
// @Description Deletes a movie by ID
// @Tags movies
// @Produce json
// @Param id path int true "Movie ID"
// @Success 200 {object} types.DeleteMovieResponse
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Security BearerAuth
// @Router /movies/{id} [delete]
func (h *MovieHandler) DeleteMovie(c *gin.Context) {
	var req types.DeleteMovieRequest
	if err := c.ShouldBindUri(&req); err != nil {
		h.log.Warn("Invalid delete movie request", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	resp, err := h.svc.DeleteMovie(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete movie"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
