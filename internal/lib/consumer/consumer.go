package consumer

import (
	"encoding/json"

	amqplib "github.com/rabbitmq/amqp091-go"
	cons "github.com/sagarmaheshwary/microservices-video-catalog-service/internal/constant"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/broker"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/log"
)

var C *Consumer

type Consumer struct {
	channel *amqplib.Channel
}

func (c *Consumer) Consume() {
	q, err := c.channel.QueueDeclare(
		cons.QueueVideoCatalogService, // name
		true,                          // durable
		false,                         // delete when unused
		false,                         // exclusive
		false,                         // no-wait
		nil,                           // arguments
	)

	if err != nil {
		log.Error("AMQP queue error %v", err)
	}

	messages, err := c.channel.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	if err != nil {
		log.Fatal("Broker queue listen failed %v", err)
	}

	log.Info("Broker listening on queue %q", cons.QueueVideoCatalogService)

	var forever chan struct{}

	go func() {
		for message := range messages {
			s := broker.MessageType{}
			json.Unmarshal(message.Body, &s)
			log.Info("Message received %#v", s.Data)

			// switch s.Key {
			// case:

			// }
		}
	}()

	<-forever
}

func Init(channel *amqplib.Channel) *Consumer {
	C = &Consumer{channel: channel}

	return C
}
