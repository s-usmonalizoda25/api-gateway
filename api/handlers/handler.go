package handlers

import (
	"github.com/s-usmonalizoda25/api-gateway/pkg/jwt"
	"github.com/s-usmonalizoda25/api-gateway/pkg/rabbitmq"
	"github.com/s-usmonalizoda25/api-gateway/services"
	"go.uber.org/zap"
)

type handler struct {
	serviceManager services.IServiceManager
	rabbitMQ       *rabbitmq.RabbitMQ
	jwtParser      *jwt.Parser
	log            *zap.Logger
}

func NewHandler(serviceManager services.IServiceManager, rabbitMQ *rabbitmq.RabbitMQ, jwtParser *jwt.Parser, log *zap.Logger) *handler {
	return &handler{
		serviceManager: serviceManager,
		rabbitMQ:       rabbitMQ,
		jwtParser:      jwtParser,
		log:            log,
	}
}
