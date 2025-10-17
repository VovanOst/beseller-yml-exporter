package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config содержит конфигурацию приложения
type Config struct {
	GraphQLEndpoint string
	ShopName        string
	ShopCompany     string
	ShopURL         string
	Currency        string
	StatusID        int
	OutputPath      string
	HTTPTimeout     time.Duration
	LogLevel        string
}

// LoadFromEnv загружает конфигурацию из переменных окружения
func LoadFromEnv() *Config {
	// Попытка загрузить .env файл (игнорируем ошибку если файл не существует)
	_ = godotenv.Load()

	cfg := &Config{
		GraphQLEndpoint: os.Getenv("GRAPHQL_ENDPOINT"),
		ShopName:        getEnvOrDefault("SHOP_NAME", "Demo Shop"),
		ShopCompany:     getEnvOrDefault("SHOP_COMPANY", "Company"),
		ShopURL:         getEnvOrDefault("SHOP_URL", "https://demo.beseller.com"),
		Currency:        getEnvOrDefault("CURRENCY", "BYN"),
		StatusID:        getEnvAsInt("STATUS_ID", 1),
		OutputPath:      getEnvOrDefault("OUTPUT_PATH", "export.yml"),
		HTTPTimeout:     getEnvAsDuration("HTTP_TIMEOUT", 30*time.Second),
		LogLevel:        getEnvOrDefault("LOG_LEVEL", "info"),
	}

	return cfg
}

// getEnvOrDefault возвращает значение переменной окружения или значение по умолчанию
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt возвращает значение переменной окружения как int или значение по умолчанию
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsDuration возвращает значение переменной окружения как duration или значение по умолчанию
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		// Пробуем распарсить как секунды
		if seconds, err := strconv.Atoi(value); err == nil {
			return time.Duration(seconds) * time.Second
		}
		// Пробуем распарсить как duration string (например "30s")
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
