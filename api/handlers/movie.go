package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/s-usmonalizoda25/api-gateway/models"
	"github.com/s-usmonalizoda25/api-gateway/pkg/errs"
	moviepb "github.com/s-usmonalizoda25/protoCinemaService/gen/movie"
)

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

func (h *handler) ListMovies(c *gin.Context) {
	response, err := h.serviceManager.MovieService().List(c.Request.Context(), &moviepb.ListMovieRequest{})
	if err != nil {
		errs.HandleError(c, h.log, errs.MsgFailedListMovies, err)
		return
	}
	c.JSON(http.StatusOK, response)
}

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
