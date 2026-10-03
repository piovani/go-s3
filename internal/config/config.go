package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// AWS configuration
	AWSRegion    string
	AWSAccessKey string
	AWSSecretKey string
	AWSEndpoint  string
	BucketName   string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: arquivo .env não encontrado. Lendo variáveis do sistema.")
	}

	return &Config{
		AWSRegion:    getEnvOrDefault("AWS_REGION", "us-east-1"),
		AWSAccessKey: os.Getenv("AWS_ACCESS_KEY_ID"),
		AWSSecretKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		AWSEndpoint:  os.Getenv("AWS_ENDPOINT_URL"),
		BucketName:   os.Getenv("AWS_BUCKET_NAME"),
	}
}

func getEnvOrDefault(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
