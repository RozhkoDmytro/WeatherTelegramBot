package main

import (
	"io"
	"log/slog"
	"os"
	"time"

	"projecttelegrambot/pkg/bot"
	"projecttelegrambot/pkg/config"

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
	telegramBot := bot.NewBot(cfg.Token, logger)

	for {

		// Get updates
		updates, err := telegramBot.GetUpdates()
		if err != nil {
			logger.Error("Failed to get updates: %v\n", err)
			return
		}

		// set UUID for this request for this child logger
		childLogger := logger.With(slog.String("UUID", uuid.New().String()))
		telegramBot.Logger = childLogger

		for _, update := range updates.Result {
			// Create and send rescponse
			msg := telegramBot.CreateResponseToCommand(update.Message.Text)
			_, err := telegramBot.Send(update.Message.Chat.ID, msg)
			if err != nil {
				childLogger.Error("Failed to send message: %v\n", err)
				return
			}
			telegramBot.Offset = update.UpdateID + 1
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
