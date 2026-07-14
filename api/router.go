package api

import (
	"github.com/gin-gonic/gin"
	"github.com/s-usmonalizoda25/api-gateway/api/handlers"
	"github.com/s-usmonalizoda25/api-gateway/config"
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
		api.GET("/user/:user_id", handler.GetUser)
		api.POST("/movie/create", handler.CreateMovie)
		api.POST("/booking/create", handler.CreateBooking)
		api.GET("/booking/:booking_id", handler.GetBooking)
		api.GET("/booking/user/:user_id", handler.GetUserBookings)
		api.DELETE("/booking/:booking_id", handler.CancelBooking)
		api.POST("/user/login", handler.Login)
	}

	return router
}
