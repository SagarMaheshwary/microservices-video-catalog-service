package server

import (
	"context"

	userrpc "github.com/sagarmaheshwary/microservices-video-catalog-service/internal/grpc/client/user"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/broker"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/database"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/logger"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/prometheus"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type healthServer struct {
	healthpb.HealthServer
}

func (h *healthServer) Check(ctx context.Context, req *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	status := getServicesHealthStatus()

	logger.Info("Overall health status: %q", status)

	response := &healthpb.HealthCheckResponse{
		Status: status,
	}

	if status == healthpb.HealthCheckResponse_NOT_SERVING {
		prometheus.ServiceHealth.Set(0)
		return response, nil
	}

	prometheus.ServiceHealth.Set(1)
	return response, nil
}

func getServicesHealthStatus() healthpb.HealthCheckResponse_ServingStatus {
	if !userrpc.HealthCheck() {
		return healthpb.HealthCheckResponse_NOT_SERVING
	}

	if !database.HealthCheck() {
		return healthpb.HealthCheckResponse_NOT_SERVING
	}

	if !broker.HealthCheck() {
		return healthpb.HealthCheckResponse_NOT_SERVING
	}

	return healthpb.HealthCheckResponse_SERVING
}
