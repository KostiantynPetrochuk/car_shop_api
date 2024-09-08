package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error to load .env: %v", err)
	}
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
