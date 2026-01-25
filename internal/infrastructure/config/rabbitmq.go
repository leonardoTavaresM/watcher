package config

import "os"

type RabbitMQConfig struct {
	URI      string
	Exchange string
	Queue    string
}

func GetRabbitMQConfig() RabbitMQConfig {
	return RabbitMQConfig{
		URI:      getEnv("RABBITMQ_URI", "amqp://guest:guest@localhost:5672/"),
		Exchange: getEnv("RABBITMQ_EXCHANGE", "file_events"),
		Queue:    getEnv("RABBITMQ_QUEUE", "file_events_queue"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
