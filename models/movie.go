package models

type CreateMovieRequest struct {
	Title       string `json:"title" binding:"required,min=1,max=100"`
	Description string `json:"description" binding:"required,max=500"`
	AgeLimit    int32  `json:"age_limit" binding:"required,gte=0,lte=21"`
	Duration    int32  `json:"duration" binding:"required,gt=0"`
}
