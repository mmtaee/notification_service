package configs

import (
	"os"
)

type Kavenegar struct {
	APIKey string
}

func loadKavenegar() Kavenegar {
	apiKey := os.Getenv("KAVENEGAR_API_KEY")
	return Kavenegar{
		APIKey: apiKey,
	}
}
