package main

import (
	"log"
	"time"

	"projecttelegrambot/pkg/bot"
	"projecttelegrambot/pkg/config"
)

const (
	defualtTimeout = 2 // in seconds
)

func main() {
	// Get confog with env
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to get correct config: %v\n", err)
	}
	// Create a new telegram bot
	telegramBot := bot.NewBot(cfg.Token)

	for {
		updates, err := telegramBot.GetUpdates()
		if err != nil {
			log.Fatalf("Failed to get updates: %v\n", err)
		}

		for _, update := range updates.Result {
			// update ChatID
			telegramBot.ChatId = int(update.Message.Chat.ID)

			// Create and send rescponse
			msg := telegramBot.CreateResponseToCommand(update.Message.Text)
			_, err := telegramBot.Send(msg)
			if err != nil {
				log.Fatalf("Failed to send message: %v\n", err)
			}
			telegramBot.Offset = update.UpdateID + 1
		}

		// Sleep for a bit before polling again
		time.Sleep(defualtTimeout * time.Second)
	}
}
