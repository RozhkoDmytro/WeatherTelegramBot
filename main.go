package main

import "projecttelegrambot/pkg/telegram"

const (
	defualtTimeout = 2 // in seconds
)

func main() {
	// create all background in one struct
	telegramApp, err := telegram.NewMyApp()
	if err != nil {
		panic(err)
	}

	telegramApp.Bot.ListenAndServe(defualtTimeout, telegramApp.SendResponse)
}
