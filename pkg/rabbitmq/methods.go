package rabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

func setChannel(conn *amqp.Connection) *amqp.Channel {
	channel, err := conn.Channel()
	if err != nil {
		panic(err)
		//	TODO: log fatal
	}
	return channel
}

func (r *RabbitMQ) SetExchange() {
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

func (r *RabbitMQ) SetQueue(queueName string) *amqp.Queue {
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

func (r *RabbitMQ) SetRoutingKeys(queue *amqp.Queue, keys []string) {
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
