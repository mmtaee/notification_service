package consumers

import (
	"context"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"notification/internal/providers"
	"notification/pkg/logger"
	"notification/pkg/rabbitmq"
	"sync"
)

type Message struct {
	Provider string   `json:"provider" validate:"required"`
	Data     struct{} `json:"data" validate:"required"`
}

func Consumer(ctx context.Context) {
	workers := ctx.Value("workers").(int)

	rbtmq := rabbitmq.Get()

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		consumerTag := uuid.New().String()
		go worker(ctx, rbtmq, &wg, consumerTag)
	}
	wg.Wait()
}

func worker(ctx context.Context, rbtmq *rabbitmq.RabbitMQ, wg *sync.WaitGroup, consumerTag string) {
	defer wg.Done()

	logger.Info("worker %s started", consumerTag)

	messages, err := rbtmq.Channel.Consume(
		rbtmq.Queue.Name,
		consumerTag, // Consumer tag
		false,       // Auto-ack
		false,       // Exclusive
		false,       // No-local
		false,       // No-wait
		nil,         // Arguments
	)
	if err != nil {
		logger.Critical("Failed to consume messages: %v", err)
	}

	validation := validator.New()

	for {
		select {
		case msg := <-messages:
			if msg.Body == nil {
				return
			}
			if rabbitmq.CheckRoutingKeys(msg.RoutingKey) {
				var data Message
				err = json.Unmarshal(msg.Body, &data)
				if err != nil {
					logger.Error("Error unmarshalling message: %v", err)
					referralMessage(msg)
					return
				}
				if err = validation.Struct(data); err != nil {
					logger.Error("Error validating message: %v", err)
					referralMessage(msg)
					return
				}
				prv := findProvider(data.Provider)
				if prv == nil {
					logger.Error("provider %v not found", data.Provider)
					referralMessage(msg)
					return
				}
				err = prv.Process(&msg)
				if err != nil {
					logger.Error("Error consuming message: %v", err)
					referralMessage(msg)
					return
				}
				err = msg.Ack(false)
				if err != nil {
					logger.Error("Error acknowledging message: %v", err)
					return
				}
				logger.Success("Consumed message: %v", data)
				return
			}
			logger.Error("routing %s key not found", msg.RoutingKey)
			referralMessage(msg)
		case <-ctx.Done():
			_ = rbtmq.Channel.Cancel(consumerTag, false) // Unregister the consumer
			return
		}
	}
}

func referralMessage(msg amqp.Delivery) {
	err := msg.Nack(false, false)
	if err != nil {
		logger.Error("Error Nack message: %v", err)
		return
	}
}

func findProvider(name string) providers.Provider {
	provider, err := providers.FindProvider(name)
	if err != nil {
		return nil
	}
	return provider
}
