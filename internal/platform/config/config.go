package config

import (
	"os"
)

type Config struct {
	Port          string
	AWSRegion     string
	DynamoDBTable string
	LogLevel      string
}

func LoadConfig() *Config {
	return &Config{
		Port:          getEnv("PORT", "8080"),
		AWSRegion:     getEnv("AWS_REGION", "us-east-1"),
		DynamoDBTable: getEnv("DYNAMODB_TABLE", "products"),
		LogLevel:      getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
