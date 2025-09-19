package util

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func MustGetEnv(key string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return ""
	}
	val := os.Getenv(key)
	if val == "" {
		panic("missing required env var: " + key)
	}
	return val
}
