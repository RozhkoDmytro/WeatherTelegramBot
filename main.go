package main

import (
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"projecttelegrambot/pkg/config"
	"projecttelegrambot/pkg/holiday"
	"projecttelegrambot/pkg/mongodb"
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

	// Create all new APIs and connection
	apiTelegram, err := telegrambot.NewBot(cfg.Token, logger)
	if err != nil {
		panic(err)
	}
	apiHoliday := holiday.NewApiHoliday(&http.Client{}, holiday.HolidayApiUrl, cfg.TokenHoliday)
	apiWeather := weather.NewApiWeather(&http.Client{}, weather.WeatherApiUrl, cfg.TokenWeather)
	apiMongoDB, err := mongodb.NewApiMongoDB(mongodb.BaseURL, logger)
	if err != nil {
		logger.Error("Can`t create connection with MongoDB", "Error", err)
	}

	// create all background in one struct
	telegramSrv := telegram.NewMyTelegramService(apiTelegram, apiHoliday, apiWeather, apiMongoDB)

	// Start ticker subscribers
	ticker := time.NewTicker(time.Hour)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				subscribers, err := apiMongoDB.GetSubsribersByTime(t.Hour())
				if err != nil {
					logger.Error("Can`t create report", "Error", err)
				}
				telegramSrv.SendReportWeather(subscribers)
			}
		}
	}()
	defer ticker.Stop()

	// Start process listnen and after-serving responce
	apiTelegram.ListenAndServe(defualtTimeout, telegramSrv.CreateSendResponse)
	done <- true
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
