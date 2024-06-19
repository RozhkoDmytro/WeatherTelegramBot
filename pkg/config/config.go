package config

import (
	"errors"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type config struct {
	Token   string `env:"TELEGRAM_BOT_TOKEN"`
	NameLog string `env:"NAME_LOG_FILE"`
}

// return token telegram bot or exit
func Load() (config, error) {
	var cfg config

	// Load .env file
	if err := godotenv.Load(); err != nil {
		return cfg, err
	}

	// Parse environment variables into Config struct
	cfg = config{}
	if err := env.Parse(&cfg); err != nil {
		return cfg, err
	}

	// Get the token from the environment variable
	if cfg.Token == "" {
		err := errors.New("bad token")
		return cfg, err
	}

	return cfg, nil
}
