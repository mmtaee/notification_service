package consumers

import (
	"context"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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
					_ = msg.Nack(false, false)
					return
				}
				if err = validation.Struct(data); err != nil {
					logger.Error("Error validating message: %v", err)
					_ = msg.Nack(false, false)
					return
				}

				var prv providers.Provider
				if prv, err = providers.FindProvider(data.Provider); err != nil {
					logger.Error("invalid provider %s", data.Provider)
					_ = msg.Nack(false, false)
					return
				} else {
					err = prv.Process(&msg)
					if err != nil {
						logger.Error("Error consuming message: %v", err)
						_ = msg.Nack(false, false)
						return
					}
				}
				_ = msg.Ack(false)
				logger.Success("Consumed message: %v", data)
				return
			}
			logger.Error("routing %s key not found", msg.RoutingKey)
			_ = msg.Nack(false, false)
		case <-ctx.Done():
			_ = rbtmq.Channel.Cancel(consumerTag, false) // Unregister the consumer
			return
		}
	}
}
