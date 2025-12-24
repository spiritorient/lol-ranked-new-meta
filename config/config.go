package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	RiotAPIKey     string
	OpenAIAPIKey   string
	ServerPort     string
	RiotAPIRegion  string
	OpenAIModel    string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	// Try to load .env file (ignore error if it doesn't exist)
	_ = godotenv.Load()

	config := &Config{
		RiotAPIKey:    getEnv("RIOT_API_KEY", ""),
		OpenAIAPIKey:  getEnv("OPENAI_API_KEY", ""),
		ServerPort:    getEnv("PORT", "8080"),
		RiotAPIRegion: getEnv("RIOT_API_REGION", "americas"), // americas, europe, asia, sea
		OpenAIModel:   getEnv("OPENAI_MODEL", "gpt-4o-mini"),
	}

	// Validate required configuration
	if config.RiotAPIKey == "" {
		return nil, fmt.Errorf("RIOT_API_KEY is required")
	}
	if config.OpenAIAPIKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY is required")
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

