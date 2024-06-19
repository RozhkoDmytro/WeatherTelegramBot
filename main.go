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

	// create logger
	child, err := createChild(cfg.NameLog)
	if err != nil {
		panic(err)
	}

	// Create a new telegram bot
	telegramBot := bot.NewBot(cfg.Token, child)

	for {
		updates, err := telegramBot.GetUpdates()
		if err != nil {
			child.Error("Failed to get updates: %v\n", err)
			return
		}

		for _, update := range updates.Result {
			// update ChatID
			telegramBot.ChatId = int(update.Message.Chat.ID)

			// Create and send rescponse
			msg := telegramBot.CreateResponseToCommand(update.Message.Text)
			_, err := telegramBot.Send(msg)
			if err != nil {
				child.Error("Failed to send message: %v\n", err)
				return
			}
			telegramBot.Offset = update.UpdateID + 1
		}

		// Sleep for a bit before polling again
		time.Sleep(defualtTimeout * time.Second)
	}
}

// Create logger and set fields
func createChild(NameLog string) (*slog.Logger, error) {
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

	child := logger.With(slog.String("UUID", uuid.New().String()))
	return child, nil
}
