package consumer

import (
	"context"
	"encoding/json"
	"time"

	amqplib "github.com/rabbitmq/amqp091-go"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/constant"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/handler"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/broker"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/logger"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
)

var C *Consumer

type Consumer struct {
	channel *amqplib.Channel
}

func (c *Consumer) Consume() error {
	q, err := c.declareQueue(constant.QueueVideoCatalogService)
	if err != nil {
		logger.Fatal("AMQP queue listen failed %v", err)
	}

	messages, err := c.channel.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Fatal("AMQP queue listen failed %v", err)
	}

	logger.Info("AMQP listening on queue %q", constant.QueueVideoCatalogService)

	var forever chan struct{}

	go func() {
		for message := range messages {
			// Extract tracing context from message headers
			carrier := make(propagation.MapCarrier)
			for k, v := range message.Headers {
				if str, ok := v.(string); ok {
					carrier[k] = str
				}
			}
			ctx := otel.GetTextMapPropagator().Extract(context.Background(), carrier)

			logger.Info("Consume headers %v", carrier)

			tracer := otel.Tracer(constant.ServiceName)
			ctx, span := tracer.Start(ctx, constant.TraceTypeRabbitMQConsume)

			m := broker.MessageType{}
			if err := json.Unmarshal(message.Body, &m); err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, "failed to unmarshal base message")
				logger.Error("Failed to unmarshal message body: %v", err)

				continue
			}

			span.SetAttributes(attribute.String("message_key", m.Key))

			logger.Info("AMQP Message received %q: %v", m.Key, m.Data)

			prometheus.TotalMessagesCounter.WithLabelValues(m.Key).Inc()
			start := time.Now()

			switch m.Key {
			case constant.MessageTypeVideoEncodingCompleted:
				type MessageType struct {
					Key  string                                `json:"key"`
					Data handler.VideoEncodingCompletedMessage `json:"data"`
				}
				d := new(MessageType)

				if err := json.Unmarshal(message.Body, d); err != nil {
					span.RecordError(err)
					span.SetStatus(codes.Error, "failed to unmarshal message")
					logger.Error("Failed to unmarshal message %s: %s", m.Key, err)
					continue
				}

				err := handler.ProcessVideoEncodingCompletedMessage(ctx, &d.Data)
				if err == nil {
					message.Ack(false)
					prometheus.MessageProcessingDuration.WithLabelValues(m.Key).Observe(time.Since(start).Seconds())
					span.SetStatus(codes.Ok, "message processed successfully")
				} else {
					span.RecordError(err)
					span.SetStatus(codes.Error, "message processing failed")
					logger.Error("Failed to process message %s: %s", m.Key, err)
					prometheus.MessageProcessingErrorsCounter.WithLabelValues(m.Key, err.Error()).Inc()
				}
			default:
				span.AddEvent("unknown message type")
				logger.Warn("Unknown message key: %s", m.Key)
			}

			span.End()
		}
	}()

	<-forever
	return nil
}

func (c *Consumer) declareQueue(queue string) (*amqplib.Queue, error) {
	q, err := c.channel.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		logger.Error("AMQP declare queue error %v", err)

		return nil, err
	}

	return &q, err
}

func Init(channel *amqplib.Channel) *Consumer {
	C = &Consumer{channel: channel}

	return C
}
