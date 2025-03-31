package configs

type Config struct {
	RabbitMQ RabbitMQConfig
}

type Kavenegar struct {
	APIKey string
}

type YCloud struct {
	APIKey string
}

var config Config

func LoadConfig() {
	config = Config{
		RabbitMQ: loadRabbitMQ(),
	}
}

func GetRabbitMQConfig() RabbitMQConfig {
	return config.RabbitMQ
}
