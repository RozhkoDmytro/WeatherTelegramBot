package config

import (
	"log"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type config struct {
	Token string `env:"TELEGRAM_BOT_TOKEN"`
}

// return token telegram bot or exit
func MustToken() string {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Parse environment variables into Config struct
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Error parsing environment variables: %v", err)
	}

	// Get the token from the environment variable
	if cfg.Token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable not set")
	}

	return cfg.Token
}
