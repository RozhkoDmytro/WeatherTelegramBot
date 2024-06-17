package main

import (
	"fmt"
	"gitlab/foxstudent107249/projecttelegrambot/pkg/bot"
	"log"
	"time"

	"github.com/joho/godotenv"
)

const (
	telegramTokenEnv = "TELEGRAM_BOT_TOKEN"
	defualtTimeout   = 2 // in seconds
)

var (
	telegramToken = ""
	myEnv         map[string]string
)

// Set all parametrs
func init() {
	// Set token
	telegramToken = mustToken()
}

func main() {
	offset := 0

	for {
		updates, err := bot.GetUpdates(telegramToken, offset)
		if err != nil {
			fmt.Printf("Failed to get updates: %v\n", err)
			return
		}

		for _, update := range updates.Result {
			msg := bot.RunCommand(update)
			bot.Send(int(update.Message.Chat.ID), msg, telegramToken)
			offset = update.UpdateID + 1
		}

		// Sleep for a bit before polling again
		time.Sleep(defualtTimeout * time.Second)
	}
}

// return token telegram bot or exit
func mustToken() string {
	// Initialized ENV from file
	var err error
	myEnv, err = godotenv.Read()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get the token from the environment variable
	token := myEnv[telegramTokenEnv]
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable not set")
	}

	return token
}
