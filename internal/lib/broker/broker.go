package broker

import (
	"context"
	"fmt"
	"sync"
	"time"

	amqplib "github.com/rabbitmq/amqp091-go"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/config"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/consumer"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/logger"
)

var Conn *amqplib.Connection

var (
	reconnectLock sync.Mutex
	retries       = 10
	interval      = 5
)

func MaintainConnection(ctx context.Context) {
	if err := connect(); err != nil {
		logger.Error("Initial AMQP connection attempt failed: %v", err)
	}

	t := time.NewTicker(time.Second * time.Duration(interval))
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			if err := tryReconnect(); err != nil {
				return
			}
		}
	}
}

func connect() error {
	c := config.Conf.AMQP

	address := fmt.Sprintf("amqp://%s:%s@%s:%d", c.Username, c.Password, c.Host, c.Port)

	var err error

	Conn, err = amqplib.Dial(address)

	if err != nil {
		logger.Error("Broker connection error %v", err)

		return err
	}

	logger.Info("Broker connected on %q", address)

	channel, err := NewChannel()

	if err != nil {
		logger.Error("Unable to create listen channel %v", err)

		return err
	}

	go func() {
		consumer.Init(channel).Consume()
	}()

	return nil
}

func tryReconnect() error {
	reconnectLock.Lock()
	defer reconnectLock.Unlock()

	if HealthCheck() {
		return nil
	}

	for i := 0; i < retries; i++ {
		logger.Info("AMQP connection attempt: %d", i+1)

		if err := connect(); err == nil {
			return nil
		}

		time.Sleep(time.Duration(interval*(i+1)) * time.Second)
	}

	return fmt.Errorf("could not reconnect after %d retries", retries)
}

func NewChannel() (*amqplib.Channel, error) {
	c, err := Conn.Channel()

	if err != nil {
		logger.Error("Broker channel error %v", err)

		return nil, err
	}

	return c, nil
}

func HealthCheck() bool {
	if Conn == nil || Conn.IsClosed() {
		logger.Warn("AMQP health check failed!")

		return false
	}

	return true
}
