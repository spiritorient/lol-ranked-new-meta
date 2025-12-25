package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	RiotAPIKey           string
	OpenAIAPIKey         string
	ServerPort           string
	RiotAPIRegion        string
	OpenAIModel          string
	AnalyticsDataPath    string
	AnalyticsMaxDays     int  // Maximum days to keep requests (0 = unlimited)
	AnalyticsMaxRecords  int  // Maximum total records to keep (0 = unlimited)
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	// Try to load .env file (ignore error if it doesn't exist)
	_ = godotenv.Load()

	config := &Config{
		RiotAPIKey:        getEnv("RIOT_API_KEY", ""),
		OpenAIAPIKey:      getEnv("OPENAI_API_KEY", ""),
		ServerPort:        getEnv("PORT", "8080"),
		RiotAPIRegion:     getEnv("RIOT_API_REGION", "americas"), // americas, europe, asia, sea
		OpenAIModel:       getEnv("OPENAI_MODEL", "gpt-4o-mini"),
		// Default to /data/analytics.json for Render.com persistent disk
		// For local development, use ./data/analytics.json
		AnalyticsDataPath:   getEnv("ANALYTICS_DATA_PATH", "/data/analytics.json"),
		AnalyticsMaxDays:    getEnvInt("ANALYTICS_MAX_DAYS", 0),      // 0 = unlimited
		AnalyticsMaxRecords: getEnvInt("ANALYTICS_MAX_RECORDS", 0),   // 0 = unlimited
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

func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

