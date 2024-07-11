package config

import (
	"errors"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Config struct {
	Token        string `env:"TELEGRAM_BOT_TOKEN"`
	NameLog      string `env:"NAME_LOG_FILE"`
	TokenHoliday string `env:"HOLIDAY_TOKEN"`
	TokenWeather string `env:"WEATHER_TOKEN"`
}

// return token telegram bot or exit
func Load() (Config, error) {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		return Config{}, err
	}

	// Parse environment variables into Config struct
	cfg := Config{}
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
