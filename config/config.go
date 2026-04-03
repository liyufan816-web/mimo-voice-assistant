package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MimoAPIKey   string
	MimoEndpoint string
	ServerPort   string
	MaxFileSize  int64
}

func LoadConfig() (*Config, error) {
	// 加载.env文件（开发环境）
	godotenv.Load()

	return &Config{
		MimoAPIKey:   getEnv("MIMO_API_KEY", ""),
		MimoEndpoint: getEnv("MIMO_ENDPOINT", "https://api.mimo.ai/v2/tts"),
		ServerPort:   getEnv("SERVER_PORT", "8080"),
		MaxFileSize:  10 * 1024 * 1024, // 10MB
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
