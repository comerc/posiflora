package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	DB       DBConfig
	Server   ServerConfig
	Telegram TelegramConfig
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
}

type ServerConfig struct {
	Host        string
	Port        int
	CORSOrigins []string
}

type TelegramConfig struct {
	MockMode bool
}

func Load() (*Config, error) {
	cfg := &Config{}

	// Database config
	cfg.DB.Host = getEnv("DB_HOST", "localhost")
	cfg.DB.Port = getEnvAsInt("DB_PORT", 5432)
	cfg.DB.User = getEnv("DB_USER", "postgres")
	cfg.DB.Password = getEnv("DB_PASSWORD", "posiflora")
	cfg.DB.Database = getEnv("DB_NAME", "posiflora")
	cfg.DB.SSLMode = getEnv("DB_SSLMODE", "disable")

	// Server config
	cfg.Server.Host = getEnv("SERVER_HOST", "0.0.0.0")
	cfg.Server.Port = getEnvAsInt("SERVER_PORT", 8080)
	cfg.Server.CORSOrigins = getEnvAsStringSlice("CORS_ORIGINS", []string{"http://localhost:5173"})

	// Telegram config
	cfg.Telegram.MockMode = getEnvAsBool("TELEGRAM_MOCK_MODE", true)

	return cfg, nil
}

func (c *DBConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Database, c.SSLMode)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsStringSlice(key string, defaultValue []string) []string {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	// Разделяем по запятой и очищаем пробелы
	values := []string{}
	for _, v := range strings.Split(valueStr, ",") {
		trimmed := strings.TrimSpace(v)
		if trimmed != "" {
			values = append(values, trimmed)
		}
	}
	if len(values) == 0 {
		return defaultValue
	}
	return values
}
