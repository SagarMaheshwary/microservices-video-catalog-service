package main

import (
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/config"
	userrpc "github.com/sagarmaheshwary/microservices-video-catalog-service/internal/grpc/client/user"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/grpc/server"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/broker"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/consumer"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/database"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/logger"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/prometheus"
)

func main() {
	logger.Init()
	config.Init()

	database.Connect()

	broker.Connect()
	defer broker.Conn.Close()

	channel, err := broker.NewChannel()

	if err != nil {
		logger.Fatal("Unable to create listen channel %v", err)
	}

	c := consumer.Init(channel)

	go func() {
		c.Consume()
	}()

	userrpc.Connect()

	go func() {
		prometheus.Connect()
	}()

	server.Connect()
}
