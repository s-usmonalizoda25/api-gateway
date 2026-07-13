package errs

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

var (
	MsgFailedGetUser         = "failed to get user"
	MsgFailedCreateMovie     = "failed to create movie"
	MsgFailedCreateBooking   = "failed to create booking"
	MsgFailedGetBooking      = "failed to get booking"
	MsgFailedGetUserBookings = "failed to get user bookings"
	MsgFailedCancelBooking   = "failed to cancel booking"
)

func HandleValidationError(c *gin.Context, err error) {
	if _, ok := err.(validator.ValidationErrors); ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
}

func HandleError(c *gin.Context, log *zap.Logger, msg string, err error) {
	log.Error(msg, zap.Error(err))
	c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
}
