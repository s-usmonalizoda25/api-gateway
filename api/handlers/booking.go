package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/s-usmonalizoda25/api-gateway/models"
	"github.com/s-usmonalizoda25/api-gateway/pkg/errs"
	bookingpb "github.com/s-usmonalizoda25/protoCinemaService/gen/booking"
)

func (h *handler) CreateBooking(c *gin.Context) {
	var body models.CreateBookingRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		errs.HandleValidationError(c, err)
		return
	}

	response, err := h.serviceManager.BookingService().CreateBooking(c.Request.Context(), &bookingpb.CreateBookingRequest{
		UserId: body.UserID, MovieId: body.MovieID,
	})

	if err != nil {
		errs.HandleError(c, h.log, errs.MsgFailedCreateBooking, err)
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *handler) GetBooking(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("booking_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking_id"})
		return
	}
	response, err := h.serviceManager.BookingService().GetBooking(c.Request.Context(), &bookingpb.GetBookingRequest{Id: id})
	if err != nil {
		errs.HandleError(c, h.log, errs.MsgFailedGetBooking, err)
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h *handler) GetUserBookings(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}
	response, err := h.serviceManager.BookingService().GetUserBookings(c.Request.Context(), &bookingpb.GetUserBookingsRequest{UserId: id})
	if err != nil {
		errs.HandleError(c, h.log, errs.MsgFailedGetUserBookings, err)
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h *handler) CancelBooking(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("booking_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking_id"})
		return
	}
	_, err = h.serviceManager.BookingService().CancelBooking(c.Request.Context(), &bookingpb.CancelBookingRequest{Id: id})
	if err != nil {
		errs.HandleError(c, h.log, errs.MsgFailedCancelBooking, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "booking cancelled successfully"})
}

func (h *handler) GetMyBookings(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id not found in token"})
		return
	}

	uid := int64(userID.(float64))

	response, err := h.serviceManager.BookingService().GetUserBookings(
		c.Request.Context(),
		&bookingpb.GetUserBookingsRequest{UserId: uid},
	)
	if err != nil {
		errs.HandleError(c, h.log, errs.MsgFailedGetUserBookings, err)
		return
	}
	c.JSON(http.StatusOK, response)
}
