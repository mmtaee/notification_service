package configs

import (
	"os"
)

type RabbitMQConfig struct {
	Host     string
	Port     string
	Scheme   string
	Username string
	Password string
	Vhost    string
	Exchange string
	Queue    string
}

func loadRabbitMQ() RabbitMQConfig {
	rabbitMQHost := os.Getenv("RABBITMQ_HOST")
	rabbitMQPort := os.Getenv("RABBITMQ_PORT")
	rabbitMQSecure := os.Getenv("RABBITMQ_SECURE")
	rabbitMQUser := os.Getenv("RABBITMQ_USER")
	rabbitMQPassword := os.Getenv("RABBITMQ_PASSWORD")
	rabbitMQVhost := os.Getenv("RABBITMQ_VHOST")
	rabbitMQExchange := os.Getenv("RABBITMQ_EXCHANGE")
	rabbitMQQueue := os.Getenv("RABBITMQ_QUEUE")

	if rabbitMQHost == "" || rabbitMQUser == "" || rabbitMQPassword == "" || rabbitMQExchange == "" || rabbitMQQueue == "" {
		//	TODO: add LogFatal
	}

	if rabbitMQPort == "" {
		rabbitMQPort = "5672"
	}

	rabbitMQScheme := "amqps"
	if rabbitMQSecure == "" || rabbitMQSecure == "0" {
		rabbitMQScheme = "amqp"
	}

	if rabbitMQVhost == "" {
		rabbitMQVhost = "/"
	}

	return RabbitMQConfig{
		Host:     rabbitMQHost,
		Port:     rabbitMQPort,
		Scheme:   rabbitMQScheme,
		Username: rabbitMQUser,
		Password: rabbitMQPassword,
		Vhost:    rabbitMQVhost,
		Exchange: rabbitMQExchange,
		Queue:    rabbitMQQueue,
	}
}
