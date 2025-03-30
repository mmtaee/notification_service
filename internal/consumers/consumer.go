package consumers

import (
	"context"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"notification/internal/handlers"
	"notification/pkg/rabbitmq"
	"slices"
	"sync"
)

func Consumer(ctx context.Context, rbtmq *rabbitmq.RabbitMQ) {
	var wg sync.WaitGroup

	queueName := ctx.Value("queueName").(string)
	workers := ctx.Value("workers").(int)

	rbtmq.SetExchange()
	queue := rbtmq.SetQueue(queueName)
	rbtmq.SetRoutingKeys(queue, handlers.GetRoutingKeys())

	for i := 0; i < workers; i++ {
		wg.Add(1)
		consumerTag := uuid.New().String()
		go worker(ctx, rbtmq, queue, &wg, consumerTag)
	}
	wg.Wait()
}

func worker(ctx context.Context, rbtmq *rabbitmq.RabbitMQ, queue *amqp.Queue, wg *sync.WaitGroup, consumerTag string) {
	defer wg.Done()

	log.Println("worker", consumerTag, "started")

	handler := handlers.NewHandler()

	messages, err := rbtmq.Channel.Consume(
		queue.Name,
		consumerTag, // Consumer tag
		false,       // Auto-ack
		false,       // Exclusive
		false,       // No-local
		false,       // No-wait
		nil,         // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}

	for {
		select {
		case msg := <-messages:
			if msg.Body == nil {
				return
			}
			log.Printf("Received a message: %s", msg.Body)
			if slices.Contains(handlers.GetRoutingKeys(), msg.RoutingKey) {
				handler.Send(&msg)
			} else {
				log.Println("routing key not found")
				//	TODO: add log warning
				msg.Nack(false, false)
			}
		case <-ctx.Done():
			log.Println("🛑 Context canceled, stopping consumer...")
			_ = rbtmq.Channel.Cancel(consumerTag, false) // Unregister the consumer
			return
		}
	}
}
