package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/s-usmonalizoda25/api-gateway/models"
	"github.com/s-usmonalizoda25/api-gateway/pkg/errs"
	bookingpb "github.com/s-usmonalizoda25/protoCinemaService/gen/booking"
)

func getAuthContext(c *gin.Context) (userID int64, isAdmin bool, ok bool) {
	uidRaw, exists := c.Get("user_id")
	if !exists {
		return 0, false, false
	}
	uidFloat, valid := uidRaw.(float64)
	if !valid {
		return 0, false, false
	}

	roleRaw, _ := c.Get("role")
	admin := false
	if roleFloat, valid := roleRaw.(float64); valid {
		admin = int(roleFloat) == 2
	}

	return int64(uidFloat), admin, true
}

// CreateBooking
// @Summary Create a booking
// @Description Book a movie ticket
// @Tags Booking
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.CreateBookingRequest true "Booking Info"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/booking/create [post]
func (h *handler) CreateBooking(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		errs.HandleAuthError(c, h.log, errs.MsgUnauthorized)
		return
	}
	uidFloat, ok := userID.(float64)
	if !ok {
		errs.HandleAuthError(c, h.log, errs.MsgUnauthorized)
		return
	}
	uid := int64(uidFloat)

	var body models.CreateBookingRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		errs.HandleValidationError(c, err)
		return
	}

	response, err := h.serviceManager.BookingService().CreateBooking(c.Request.Context(), &bookingpb.CreateBookingRequest{
		UserId: uid, MovieId: body.MovieID,
	})
	if err != nil {
		errs.HandleError(c, h.log, errs.MsgFailedCreateBooking, err)
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetBooking
// @Summary Get Booking by ID
// @Description Get booking details by its ID
// @Tags Booking
// @Security BearerAuth
// @Produce json
// @Param booking_id path int true "Booking ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/booking/{booking_id} [get]

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

// GetUserBookings
// @Summary Get User Bookings
// @Description Get all bookings for a specific user ID
// @Tags Booking
// @Security BearerAuth
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/booking/user/{user_id} [get]
func (h *handler) GetUserBookings(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	callerID, isAdmin, ok := getAuthContext(c)
	if !ok {
		errs.HandleAuthError(c, h.log, errs.MsgUnauthorized)
		return
	}
	if !isAdmin && callerID != id {
		errs.HandleForbiddenError(c, h.log, errs.MsgForbidden)
		return
	}

	response, err := h.serviceManager.BookingService().GetUserBookings(c.Request.Context(), &bookingpb.GetUserBookingsRequest{UserId: id})
	if err != nil {
		errs.HandleError(c, h.log, errs.MsgFailedGetUserBookings, err)
		return
	}
	c.JSON(http.StatusOK, response)
}

// CancelBooking
// @Summary Cancel Booking
// @Description Cancel a booking by ID
// @Tags Booking
// @Security BearerAuth
// @Produce json
// @Param booking_id path int true "Booking ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/booking/{booking_id} [delete]
func (h *handler) CancelBooking(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("booking_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking_id"})
		return
	}

	callerID, isAdmin, ok := getAuthContext(c)
	if !ok {
		errs.HandleAuthError(c, h.log, errs.MsgUnauthorized)
		return
	}

	existing, err := h.serviceManager.BookingService().GetBooking(c.Request.Context(), &bookingpb.GetBookingRequest{Id: id})
	if err != nil {
		errs.HandleError(c, h.log, errs.MsgFailedGetBooking, err)
		return
	}
	if !isAdmin && existing.Booking.UserId != callerID {
		errs.HandleForbiddenError(c, h.log, errs.MsgForbidden)
		return
	}

	_, err = h.serviceManager.BookingService().CancelBooking(c.Request.Context(), &bookingpb.CancelBookingRequest{Id: id})
	if err != nil {
		errs.HandleError(c, h.log, errs.MsgFailedCancelBooking, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "booking cancelled successfully"})
}

// GetMyBookings
// @Summary Get My Bookings
// @Description Get bookings for the currently authenticated user
// @Tags Booking
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/booking/me [get]
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
