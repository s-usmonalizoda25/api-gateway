package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/s-usmonalizoda25/api-gateway/api"
	"github.com/s-usmonalizoda25/api-gateway/config"
	"github.com/s-usmonalizoda25/api-gateway/pkg/logger"
	"github.com/s-usmonalizoda25/api-gateway/services"
	"go.uber.org/zap"
)

func main() {
	conf, err := config.New("./config/config.env")
	if err != nil {
		log.Fatal(err)
	}

	myLogger := logger.New()
	defer myLogger.Sync()

	serviceManager, err := services.NewServiceManager(conf.Services)
	if err != nil {
		myLogger.Fatal("services.NewServiceManager():", zap.Error(err))
	}

	handler := api.New(api.Option{
		ServiceManager: serviceManager,
		Log:            myLogger.Logger,
	})

	srv := &http.Server{
		Addr:    conf.HTTPPORT,
		Handler: handler,
	}

	go func() {
		myLogger.Info("server started", zap.String("port", conf.HTTPPORT))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			myLogger.Fatal("failed to listen", zap.Error(err))
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	myLogger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		myLogger.Fatal("server forced to shutdown", zap.Error(err))
	}

	myLogger.Info("server stopped")
}
