package user

import (
	"context"

	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/config"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/log"
	pb "github.com/sagarmaheshwary/microservices-video-catalog-service/internal/proto/user"
)

var User *userClient

type userClient struct {
	client pb.UserServiceClient
}

func (u *userClient) FindById(data *pb.FindByIdRequest) (*pb.FindByIdResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Conf.GRPCClient.Timeout)

	defer cancel()

	response, err := u.client.FindById(ctx, data)

	if err != nil {
		log.Error("gRPC userClient.FindById request failed: %v", err)
		return nil, err
	}

	log.Info("gRPC userClient.FindById response: %v", response)

	return response, nil
}
