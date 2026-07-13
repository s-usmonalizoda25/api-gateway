package handlers

import (
	"net/http"

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
