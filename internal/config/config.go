package config

import (
	"os"
	"strconv"
)

type Config struct {
	DBUrl string
	Port  string
}

func Load() Config {
	return Config{
		DBUrl: getEnv("DB_URL", ""),
		Port:  getEnv("PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func getEnvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}
