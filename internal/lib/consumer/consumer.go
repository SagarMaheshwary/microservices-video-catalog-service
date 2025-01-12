package consumer

import (
	"encoding/json"

	amqplib "github.com/rabbitmq/amqp091-go"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/constant"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/handler"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/broker"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/logger"
)

var C *Consumer

type Consumer struct {
	channel *amqplib.Channel
}

func (c *Consumer) Consume() {
	queue, err := c.channel.QueueDeclare(
		constant.QueueVideoCatalogService, // name
		true,                              // durable
		false,                             // delete when unused
		false,                             // exclusive
		false,                             // no-wait
		nil,                               // arguments
	)

	if err != nil {
		logger.Error("AMQP queue error %v", err)
	}

	messages, err := c.channel.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)

	if err != nil {
		logger.Fatal("Broker queue listen failed %v", err)
	}

	logger.Info("Broker listening on queue %q", constant.QueueVideoCatalogService)

	var forever chan struct{}

	go func() {
		for message := range messages {
			m := broker.MessageType{}

			json.Unmarshal(message.Body, &m)

			logger.Info("AMQP message received json %q: %s", m.Key, message.Body)

			switch m.Key {
			case constant.MessageTypeVideoEncodingCompleted:
				type MessageType struct {
					Key  string                                `json:"key"`
					Data handler.VideoEncodingCompletedMessage `json:"data"`
				}

				d := new(MessageType)

				json.Unmarshal(message.Body, &d)

				err := handler.ProcessVideoEncodingCompletedMessage(&d.Data)

				if err == nil {
					message.Ack(false)
				} else {
					logger.Error("Failed to process message %s: %s", m.Key, err)
				}
			}
		}
	}()

	<-forever
}

func Init(channel *amqplib.Channel) *Consumer {
	C = &Consumer{channel: channel}

	return C
}
