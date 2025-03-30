package rabbitmq

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"notification/pkg/configs"
)

type RabbitMQ struct {
	conn     *amqp.Connection
	Channel  *amqp.Channel
	exchange string
}

func NewRabbitMQ() *RabbitMQ {
	cfg := configs.GetRabbitMQConfig()
	conn := connection(cfg)
	return &RabbitMQ{
		conn:     conn,
		Channel:  setChannel(conn),
		exchange: cfg.Exchange,
	}
}

func connection(cfg configs.RabbitMQConfig) *amqp.Connection {
	c, err := amqp.Dial(
		fmt.Sprintf(
			"%s://%s:%s@%s:%s%s",
			cfg.Scheme, cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Vhost,
		),
	)
	if err != nil {
		log.Fatal(err)
		//	TODO: log Fatal
	}
	return c
}

func (r *RabbitMQ) GetConnection() *amqp.Connection {
	return r.conn
}

func (r *RabbitMQ) CloseConnection() {
	err := r.conn.Close()
	if err != nil {
		return
	}
}
