package main

import (
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/config"
	grpcsrv "github.com/sagarmaheshwary/microservices-video-catalog-service/internal/grpc/server"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/broker"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/consumer"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/database"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/log"
)

func main() {
	log.Init()
	config.Init()

	database.Connect()

	broker.Connect()
	defer broker.Conn.Close()

	listenChan, err := broker.NewChannel()

	if err != nil {
		log.Fatal("Unable to create listen channel %v", err)
	}

	c := consumer.Init(listenChan)

	go func() {
		c.Consume()
	}()

	grpcsrv.Connect()
}
