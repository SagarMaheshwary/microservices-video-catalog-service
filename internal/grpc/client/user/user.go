package user

import (
	"context"

	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/config"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/logger"
	userpb "github.com/sagarmaheshwary/microservices-video-catalog-service/internal/proto/user"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

var User *userClient

type userClient struct {
	client userpb.UserServiceClient
	health healthpb.HealthClient
}

func (u *userClient) FindById(data *userpb.FindByIdRequest) (*userpb.FindByIdResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Conf.GRPCClient.Timeout)

	defer cancel()

	response, err := u.client.FindById(ctx, data)

	if err != nil {
		logger.Error("gRPC userClient.FindById request failed: %v", err)
		return nil, err
	}

	logger.Info("gRPC userClient.FindById response: %v", response)

	return response, nil
}
