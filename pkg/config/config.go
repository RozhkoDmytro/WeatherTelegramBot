package config

import (
	"errors"
	"log"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type config struct {
	Token string `env:"TELEGRAM_BOT_TOKEN"`
}

// return token telegram bot or exit
func Load() (config, error) {
	var cfg config

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
		return cfg, err
	}

	// Parse environment variables into Config struct
	cfg = config{}
	if err := env.Parse(&cfg); err != nil {
		log.Printf("Error parsing environment variables: %v", err)
		return cfg, err
	}

	// Get the token from the environment variable
	if cfg.Token == "" {
		err := errors.New("TELEGRAM_BOT_TOKEN environment variable not set")
		log.Printf("Check token: %v", err)
		return cfg, err
	}

	return cfg, nil
}
