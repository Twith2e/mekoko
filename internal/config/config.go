package config

import (
	"os"
	"strconv"
)

type Config struct {
	DBUrl         string
	Port          string
	AccessSecret  string
	RefreshSecret string
	IsProd        string
}

func Load() Config {
	return Config{
		DBUrl:         getEnv("DB_URL", ""),
		Port:          getEnv("PORT", "8080"),
		AccessSecret:  getEnv("ACCESS_SECRET", ""),
		RefreshSecret: getEnv("REFRESH_SECRET", ""),
		IsProd:        getEnv("IS_PROD", "false"),
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
