package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/config"
	userrpc "github.com/sagarmaheshwary/microservices-video-catalog-service/internal/grpc/client/user"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/grpc/server"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/broker"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/database"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/jaeger"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/logger"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/prometheus"
)

func main() {
	logger.Init()
	config.Init()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	shutdownJaeger := jaeger.Init(ctx)

	database.Connect()

	promServer := prometheus.NewServer()
	go prometheus.Serve(promServer)

	go broker.MaintainConnection(ctx)

	userrpc.NewClient(ctx)

	grpcServer := server.NewServer()
	go server.Serve(grpcServer)

	<-ctx.Done()

	logger.Info("Shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := promServer.Shutdown(shutdownCtx); err != nil {
		logger.Warn("Prometheus server shutdown error: %v", err)
	}

	shutdownCtx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := shutdownJaeger(shutdownCtx); err != nil {
		logger.Warn("failed to shutdown jaeger tracer: %v", err)
	}

	grpcServer.GracefulStop()

	database.Close()

	logger.Info("Shutdown complete")
}
