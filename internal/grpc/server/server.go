package server

import (
	"fmt"
	"net"

	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/config"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/log"
	pb "github.com/sagarmaheshwary/microservices-video-catalog-service/internal/proto/video_catalog"
	"google.golang.org/grpc"
)

func Connect() {
	c := config.GetgrpcServer()

	address := fmt.Sprintf("%s:%d", c.Host, c.Port)

	listener, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatal("Failed to create tcp listner on %q: %v", address, err)
	}

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)

	pb.RegisterVideoCatalogServiceServer(grpcServer, &videoCatalogServer{})

	log.Info("gRPC server started on %q", address)

	if err := grpcServer.Serve(listener); err != nil {
		log.Error("gRPC server failed to start %v", err)
	}
}
