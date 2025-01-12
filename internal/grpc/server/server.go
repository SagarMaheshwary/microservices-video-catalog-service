package server

import (
	"fmt"
	"net"

	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/config"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/logger"
	videocatalogpb "github.com/sagarmaheshwary/microservices-video-catalog-service/internal/proto/video_catalog"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func Connect() {
	c := config.Conf.GRPCServer

	address := fmt.Sprintf("%s:%d", c.Host, c.Port)

	listener, err := net.Listen("tcp", address)

	if err != nil {
		logger.Fatal("Failed to create tcp listner on %q: %v", address, err)
	}

	var options []grpc.ServerOption

	server := grpc.NewServer(options...)

	videocatalogpb.RegisterVideoCatalogServiceServer(server, &videoCatalogServer{})
	healthpb.RegisterHealthServer(server, &healthServer{})

	logger.Info("gRPC server started on %q", address)

	if err := server.Serve(listener); err != nil {
		logger.Error("gRPC server failed to start %v", err)
	}
}
