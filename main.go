package main

import (
	"io"
	"log/slog"
	"net/http"
	"os"

	"projecttelegrambot/pkg/config"
	"projecttelegrambot/pkg/holiday"
	"projecttelegrambot/pkg/telegram"
	"projecttelegrambot/pkg/weather"

	"git.foxminded.ua/foxstudent107249/telegrambot"
)

const (
	defualtTimeout = 2 // in seconds
)

func main() {
	// Get config with env
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	// Create logger
	logger, err := createLogger(cfg.NameLog)
	if err != nil {
		panic(err)
	}

	// Create a new telegram bot
	bot, err := telegrambot.NewBot(cfg.Token, logger)
	if err != nil {
		panic(err)
	}
	apiHoliday := holiday.NewApiHoliday(&http.Client{}, holiday.HolidayApiUrl, cfg.TokenHoliday)
	apiWeather := weather.NewApiWeather(&http.Client{}, weather.WeatherApiUrl, cfg.TokenWeather)

	// create all background in one struct
	telegramApp := telegram.NewMyTelegramApp(&cfg, bot, apiHoliday, apiWeather)
	bot.ListenAndServe(defualtTimeout, telegramApp.SendResponse)
}

// Create logger and set fields
func createLogger(NameLog string) (*slog.Logger, error) {
	// Create logger
	file, err := os.OpenFile(NameLog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		return nil, err
	}

	w := io.MultiWriter(os.Stderr, file)
	handler := slog.NewJSONHandler(w, &slog.HandlerOptions{
		AddSource: true,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger, nil
}
