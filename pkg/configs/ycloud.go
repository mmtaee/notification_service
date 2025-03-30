package configs

import "os"

type YCloud struct {
	APIKey string
}

func loadYCloud() YCloud {
	apiKey := os.Getenv("YCLOUD_API_KEY")
	return YCloud{
		APIKey: apiKey,
	}
}
