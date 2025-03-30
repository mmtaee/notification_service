package configs

import "log"

type Config struct {
	RabbitMQ  RabbitMQConfig
	Providers Providers
}

type Providers struct {
	Kavenegar Kavenegar
	YCloud    YCloud
}

var config Config

func LoadConfig() {
	config = Config{
		RabbitMQ: loadRabbitMQ(),
		Providers: Providers{
			Kavenegar: loadKavenegar(),
			YCloud:    loadYCloud(),
		},
	}
}

func GetConfig() Config {
	return config
}

func GetRabbitMQConfig() RabbitMQConfig {
	return config.RabbitMQ
}

func KavenegarApiKey() string {
	apiKey := config.Providers.Kavenegar.APIKey
	if apiKey == "" {
		log.Fatal("KAVENEGAR_API_KEY environment variable not set")
	}
	return apiKey
}

func YcloudApiKey() string {
	apiKey := config.Providers.YCloud.APIKey
	if apiKey == "" {
		log.Fatal("YCloud_API_KEY environment variable not set")
	}
	return apiKey
}
