package models

type CreateBookingRequest struct {
	UserID  int64 `json:"user_id" binding:"required,gt=0"`
	MovieID int64 `json:"movie_id" binding:"required,gt=0"`
}
