package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/s-usmonalizoda25/api-gateway/models"
	"github.com/s-usmonalizoda25/api-gateway/pkg/errs"
	moviepb "github.com/s-usmonalizoda25/protoCinemaService/gen/movie"
)

// CreateMovie
// @Summary Create a movie
// @Description Create a new movie (Admin/Protected)
// @Tags Movie
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.CreateMovieRequest true "Movie Info"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/movies [post]
func (h *handler) CreateMovie(c *gin.Context) {
	var body models.CreateMovieRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		errs.HandleValidationError(c, err)
		return
	}

	response, err := h.serviceManager.MovieService().Create(
		c.Request.Context(),
		&moviepb.CreateMovieRequest{
			Title:       body.Title,
			Description: body.Description,
			Duration:    body.Duration,
			AgeLimit:    body.AgeLimit,
		})

	if err != nil {
		errs.HandleError(c, h.log, errs.MsgFailedCreateMovie, err)
		return
	}

	c.JSON(http.StatusCreated, response)
}

// ListMovies
// @Summary List all movies
// @Description Get a list of available movies
// @Tags Movie
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/movies [get]
func (h *handler) ListMovies(c *gin.Context) {
	response, err := h.serviceManager.MovieService().List(c.Request.Context(), &moviepb.ListMovieRequest{})
	if err != nil {
		errs.HandleError(c, h.log, errs.MsgFailedListMovies, err)
		return
	}
	c.JSON(http.StatusOK, response)
}


// GetMovie
// @Summary Get movie by ID
// @Description Get detailed information about a movie
// @Tags Movie
// @Produce json
// @Param id path int true "Movie ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/movies/{id} [get]
func (h *handler) GetMovie(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie id"})
		return
	}
	response, err := h.serviceManager.MovieService().GetByID(c.Request.Context(), &moviepb.GetMovieRequest{Id: id})
	if err != nil {
		errs.HandleError(c, h.log, errs.MsgFailedGetMovie, err)
		return
	}
	c.JSON(http.StatusOK, response)
}

// UpdateMovie
// @Summary Update a movie
// @Description Update movie details by ID (Admin/Protected)
// @Tags Movie
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Movie ID"
// @Param request body models.CreateMovieRequest true "Updated Movie Info"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/movies/{id} [put]
func (h *handler) UpdateMovie(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie id"})
		return
	}
	var body models.CreateMovieRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		errs.HandleValidationError(c, err)
		return
	}

	_, err = h.serviceManager.MovieService().Update(c.Request.Context(), &moviepb.UpdateMovieRequest{
		Id:          id,
		Title:       body.Title,
		Description: body.Description,
		Duration:    body.Duration,
		AgeLimit:    body.AgeLimit,
	})
	if err != nil {
		errs.HandleError(c, h.log, errs.MsgFailedUpdateMovie, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "movie updated successfully"})
}

// DeleteMovie
// @Summary Delete a movie
// @Description Delete a movie by ID (Admin/Protected)
// @Tags Movie
// @Security BearerAuth
// @Produce json
// @Param id path int true "Movie ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/movies/{id} [delete]
func (h *handler) DeleteMovie(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie id"})
		return
	}
	_, err = h.serviceManager.MovieService().Delete(c.Request.Context(), &moviepb.DeleteMovieRequest{Id: id})
	if err != nil {
		errs.HandleError(c, h.log, errs.MsgFailedDeleteMovie, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "movie deleted successfully"})
}
