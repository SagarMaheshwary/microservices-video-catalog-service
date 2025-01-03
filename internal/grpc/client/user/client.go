package user

import (
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/config"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/log"
	pb "github.com/sagarmaheshwary/microservices-video-catalog-service/internal/proto/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Connect() {
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	address := config.Conf.GRPCClient.UserServiceURL

	conn, err := grpc.Dial(address, opts...)

	if err != nil {
		log.Error("gRPC client failed to connect on %q: %v", address, err)
	}

	log.Info("gRPC client connected on %q", address)

	User = &userClient{
		client: pb.NewUserServiceClient(conn),
	}
}
