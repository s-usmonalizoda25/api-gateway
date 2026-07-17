package errs

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

var (
	MsgFailedRegister        = "failed to register user"
	MsgFailedGetUser         = "failed to get user"
	MsgFailedCreateMovie     = "failed to create movie"
	MsgFailedCreateBooking   = "failed to create booking"
	MsgFailedGetBooking      = "failed to get booking"
	MsgFailedGetUserBookings = "failed to get user bookings"
	MsgFailedCancelBooking   = "failed to cancel booking"
	MsgFailedLogin           = "failed to login"

	MsgUnauthorized = "unauthorized: invalid or missing token"
	MsgTokenExpired = "token expired"

	MsgForbidden = "forbidden: you do not have permission"
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

func HandleAuthError(c *gin.Context, log *zap.Logger, msg string) {
	log.Warn("auth error", zap.String("message", msg))
	c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
}

func HandleForbiddenError(c *gin.Context, log *zap.Logger, msg string) {
	log.Warn("forbidden error", zap.String("message", msg))
	c.JSON(http.StatusForbidden, gin.H{"error": msg})
}
