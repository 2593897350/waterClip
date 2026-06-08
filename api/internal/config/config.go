package config

import "os"

type Config struct {
	Address          string
	ProcessorBaseURL string
}

func Load() Config {
	address := os.Getenv("API_ADDRESS")
	if address == "" {
		address = ":8080"
	}

	processorBaseURL := os.Getenv("PROCESSOR_BASE_URL")
	if processorBaseURL == "" {
		processorBaseURL = "http://localhost:8000"
	}

	return Config{
		Address:          address,
		ProcessorBaseURL: processorBaseURL,
	}
}
