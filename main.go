package main

import (
	"fmt"
	"time"

	"projecttelegrambot/pkg/bot"
)

const (
	defualtTimeout = 2 // in seconds
)

var telegramToken = ""

func main() {
	offset := 0
	telegramToken = bot.MustToken()

	for {
		updates, err := bot.GetUpdates(telegramToken, offset)
		if err != nil {
			fmt.Printf("Failed to get updates: %v\n", err)
			return
		}

		for _, update := range updates.Result {
			msg := bot.GetInfo(update.Message.Text)
			bot.Send(int(update.Message.Chat.ID), msg, telegramToken)
			offset = update.UpdateID + 1
		}

		// Sleep for a bit before polling again
		time.Sleep(defualtTimeout * time.Second)
	}
}
