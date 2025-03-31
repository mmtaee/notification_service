package providers

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"notification/internal/providers/kavenegar"
	"notification/internal/providers/ycloud"
	"notification/pkg/logger"
	"os"
	"strings"
)

type Provider interface {
	Process(msg *amqp.Delivery) error
}

var providerRegistry = map[string]Provider{}

func Register() {
	providersEnv := os.Getenv("PROVIDERS")
	if providersEnv == "" {
		logger.Critical("PROVIDERS environment variable not set")
	}

	for _, provider := range strings.Split(providersEnv, ",") {
		switch provider {
		case "kavenegar":
			providerRegistry["kavenegar"] = kavenegar.NewKavenegar()
		case "ycloud":
			providerRegistry["ycloud"] = ycloud.NewYcloud()
		default:
			logger.Critical("Unknown provider: ", provider)
		}
	}
}

func FindProvider(name string) (Provider, error) {
	provider, ok := providerRegistry[name]
	if !ok {
		return nil, fmt.Errorf("unknown provider: %s", name)
	}
	return provider, nil
}
