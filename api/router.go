package api

import (
	"github.com/gin-gonic/gin"
	"github.com/s-usmonalizoda25/api-gateway/api/handlers"
	"github.com/s-usmonalizoda25/api-gateway/config"
	"github.com/s-usmonalizoda25/api-gateway/internal/middleware"
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
	}

	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware(option.Log))
	{
		protected.GET("/user/:user_id", handler.GetUser)
		protected.POST("/movie/create", handler.CreateMovie)
		protected.POST("/booking/create", handler.CreateBooking)
		protected.GET("/booking/:booking_id", handler.GetBooking)
		protected.GET("/booking/user/:user_id", handler.GetUserBookings)
		protected.DELETE("/booking/:booking_id", handler.CancelBooking)
	}

	return router
}
