package main

import (
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"projecttelegrambot/pkg/config"
	"projecttelegrambot/pkg/holiday"
	"projecttelegrambot/pkg/telegram"

	tgbotapi "git.foxminded.ua/foxstudent107249/telegrambot"

	"github.com/google/uuid"
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
	bot, err := tgbotapi.NewBot(cfg.Token, logger)
	if err != nil {
		logger.Error("Failed to create telegram bot: %v\n", err)
		return
	}

	apiHoliday := holiday.NewApiHoliday(&http.Client{}, holiday.HolidayApiUrl, cfg.TokenHoliday)

	for {

		// Get updates
		updates, err := bot.GetUpdates()
		if err != nil {
			logger.Error("Failed to get updates: %v\n", err)
			return
		}

		// set UUID for this request for this child logger
		childLogger := logger.With(slog.String("UUID", uuid.New().String()))
		bot.Logger = childLogger

		for _, update := range updates.Result {
			// Create and send rescponse
			_, err = telegram.CreateReplayMsg(bot, apiHoliday, &update)
			if err != nil {
				childLogger.Error("Failed to send message: %v\n", err)
				return
			}
			bot.Offset = update.UpdateID + 1
		}

		// Sleep for a bit before polling again
		time.Sleep(defualtTimeout * time.Second)
	}
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
