package models

type CreateBookingRequest struct {
	MovieID int64 `json:"movie_id" binding:"required,gt=0"`
}