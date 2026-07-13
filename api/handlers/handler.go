package handlers

import (
	"github.com/s-usmonalizoda25/api-gateway/services"
	"go.uber.org/zap"
)

type handler struct {
	serviceManager services.IServiceManager
	log            *zap.Logger
}

func NewHandler(serviceManager services.IServiceManager, log *zap.Logger) *handler {
	return &handler{
		serviceManager: serviceManager,
		log:            log,
	}
}
