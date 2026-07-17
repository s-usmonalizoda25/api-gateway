package api

import (
	"github.com/gin-gonic/gin"
	"github.com/s-usmonalizoda25/api-gateway/api/handlers"
	"github.com/s-usmonalizoda25/api-gateway/config"
	"github.com/s-usmonalizoda25/api-gateway/internal/middleware"
	"github.com/s-usmonalizoda25/api-gateway/models/permission"
	"github.com/s-usmonalizoda25/api-gateway/services"
	"go.uber.org/zap"
)

type Option struct {
	Conf           config.Config
	ServiceManager services.IServiceManager
	Log            *zap.Logger
}

func New(option Option) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	handler := handlers.NewHandler(option.ServiceManager, option.Log)

	api := router.Group("/api")
	{
		api.POST("/user/register", handler.Register)
		api.POST("/user/login", handler.Login)

		api.GET("/movies", handler.ListMovies)
		api.GET("/movies/:id", handler.GetMovie)
	}

	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware(option.Log))
	{
		protected.GET("/user/:user_id", middleware.CheckPermission(option.Log, permission.UserView), handler.GetUser)
		protected.GET("/user/me", middleware.CheckPermission(option.Log, permission.UserViewMe), handler.GetMyProfile)
		protected.PUT("/user/me", middleware.CheckPermission(option.Log, permission.UserUpdate), handler.UpdateMyProfile)

		protected.POST("/movies", middleware.CheckPermission(option.Log, permission.MovieCreate), handler.CreateMovie)
		protected.PUT("/movies/:id", middleware.CheckPermission(option.Log, permission.MovieUpdate), handler.UpdateMovie)
		protected.DELETE("/movies/:id", middleware.CheckPermission(option.Log, permission.MovieDelete), handler.DeleteMovie)

		protected.POST("/booking/create", middleware.CheckPermission(option.Log, permission.BookingCreate), handler.CreateBooking)
		protected.GET("/booking/me", middleware.CheckPermission(option.Log, permission.BookingViewMe), handler.GetMyBookings)
		protected.GET("/booking/:booking_id", middleware.CheckPermission(option.Log, permission.BookingView), handler.GetBooking)
		protected.GET("/booking/user/:user_id", middleware.CheckPermission(option.Log, permission.BookingViewMe), handler.GetUserBookings)
		protected.DELETE("/booking/:booking_id", middleware.CheckPermission(option.Log, permission.BookingCancel), handler.CancelBooking)
	}

	return router
}
