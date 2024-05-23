package broker

import (
	"fmt"

	amqplib "github.com/rabbitmq/amqp091-go"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/config"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/log"
)

var Conn *amqplib.Connection

type MessageType struct {
	Key  string `json:"key"`
	Data any    `json:"data"`
}

func Connect() {
	c := config.Getamqp()

	address := fmt.Sprintf("amqp://%s:%s@%s:%d", c.Username, c.Password, c.Host, c.Port)

	var err error

	Conn, err = amqplib.Dial(address)

	if err != nil {
		log.Error("Broker connection error %v", err)
	}

	log.Info("Broker connected on %q", address)
}

func NewChannel() (*amqplib.Channel, error) {
	c, err := Conn.Channel()

	if err != nil {
		log.Error("Broker channel error %v", err)

		return nil, err
	}

	return c, nil
}
