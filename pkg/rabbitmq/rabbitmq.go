package rabbitmq

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"notification/pkg/configs"
	"notification/pkg/logger"
	"slices"
)

type RabbitMQ struct {
	conn     *amqp.Connection
	Channel  *amqp.Channel
	exchange string
	Queue    *amqp.Queue
}

var rbtmq *RabbitMQ

var routingKeys = []string{"otp.sms", "otp.call", "event.sms"}

func Init() {
	cfg := configs.GetRabbitMQConfig()
	conn := connection(cfg)
	rbtmq = &RabbitMQ{
		conn:     conn,
		Channel:  setChannel(conn),
		exchange: cfg.Exchange,
	}
	rbtmq.setExchange()
	queue := rbtmq.setQueue(cfg.Queue)
	rbtmq.Queue = queue
	rbtmq.setRoutingKeys(queue, routingKeys)
}

func Get() *RabbitMQ {
	return rbtmq
}

func connection(cfg configs.RabbitMQConfig) *amqp.Connection {
	c, err := amqp.Dial(
		fmt.Sprintf(
			"%s://%s:%s@%s:%s%s",
			cfg.Scheme, cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Vhost,
		),
	)
	if err != nil {
		logger.Critical("rabbitmq connection error:", err)
	}
	return c
}

func CloseConnection() {
	err := rbtmq.conn.Close()
	if err != nil {
		return
	}
}

func setChannel(conn *amqp.Connection) *amqp.Channel {
	channel, err := conn.Channel()
	if err != nil {
		panic(err)
		//	TODO: log fatal
	}
	return channel
}

func (r *RabbitMQ) setExchange() {
	err := r.Channel.ExchangeDeclare(
		r.exchange, // Name
		"topic",    // Type
		true,       // Durable
		false,      // Auto-deleted
		false,      // Internal
		false,      // No-wait
		nil,        // Arguments
	)
	if err != nil {
		// TODO: log fatal
		return
	}
}

func (r *RabbitMQ) setQueue(queueName string) *amqp.Queue {
	q, err := r.Channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		//	TODO: log fatal
	}
	return &q
}

func (r *RabbitMQ) setRoutingKeys(queue *amqp.Queue, keys []string) {
	for _, key := range keys {
		err := r.Channel.QueueBind(
			queue.Name,
			key,
			r.exchange,
			false,
			nil,
		)
		if err != nil {
			//	TODO: log Fatal
		}
	}
}

func CheckRoutingKeys(name string) bool {
	if slices.Contains(routingKeys, name) {
		return true
	}
	return false
}
