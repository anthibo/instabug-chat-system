package config

import (
	"os"
)

type Config struct {
	ServerPort  string
	MySQLDSN    string
	RabbitMQURL string
	RedisURL    string
}

func Load() *Config {
	return &Config{
		ServerPort:  getEnv("SERVER_PORT", "3001"),
		MySQLDSN:    getEnv("MYSQL_DSN", "root:password@tcp(localhost:3306)/chat_system"),
		RabbitMQURL: getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		RedisURL:    getEnv("REDIS_URL", "localhost:6379"),
	}
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
