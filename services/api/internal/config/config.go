package config

import (
	"os"
	"strings"
)

type Config struct {
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	Port        string
	APIKey      string
	AgentAPIKey string
}

func Load() *Config {
	return &Config{
		DBHost:      os.Getenv("DB_HOST"),
		DBPort:      os.Getenv("DB_PORT"),
		DBUser:      os.Getenv("DB_USER"),
		DBPassword:  os.Getenv("DB_PASSWORD"),
		DBName:      os.Getenv("DB_NAME"),
		Port:        envOr("PORT", "8000"),
		APIKey:      os.Getenv("API_KEY"),
		AgentAPIKey: os.Getenv("AGENT_API_KEY"),
	}
}

func envOr(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}
