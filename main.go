package main

import (
	"fmt"
	"time"

	"projecttelegrambot/config"
	"projecttelegrambot/pkg/bot"
)

const (
	defualtTimeout = 2 // in seconds
)

var telegramBot = bot.ApiTelegramBot{
	Token: config.MustToken(),
}

func main() {
	for {
		updates, err := telegramBot.GetUpdates()
		if err != nil {
			fmt.Printf("Failed to get updates: %v\n", err)
			return
		}

		for _, update := range updates.Result {
			// update ChatID
			telegramBot.ChatId = int(update.Message.Chat.ID)

			// Create and send rescponse
			msg := telegramBot.CreateResponseToCommand(update.Message.Text)
			telegramBot.Send(msg)
			telegramBot.Offset = update.UpdateID + 1
		}

		// Sleep for a bit before polling again
		time.Sleep(defualtTimeout * time.Second)
	}
}
