package main

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var (
	myEnv   map[string]string
	infoMap map[string]string
)

func init() {
	helpStartInfo := `
/start   - get information about all bot commands
/help    - too same like start
/about  - get some information about me
/links   - send my(developer) links`

	linksInfo := `
https://www.linkedin.com/in/dmytro-rozhko-bas-1c-golang-junior/
https://animated-panda-0382af.netlify.app/
	`
	infoMap = make(map[string]string)
	infoMap["start"] = helpStartInfo
	infoMap["help"] = helpStartInfo
	infoMap["about"] = "Rozhko Dmytro; Go developer; bad character; unmarried (C)"
	infoMap["links"] = linksInfo
}

func main() {
	// Create bot
	t := mustToken()
	bot := createBot(t)

	// fixed info about name telegram bot
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// main logic: w8 for command and send answer
	u := tgbotapi.NewUpdate(60)
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			if update.Message.IsCommand() {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, infoMap[update.Message.Command()])
				bot.Send(msg)
			}
		}
	}
}

/* func sendAnswer(update tgbotapi.Update) {
} */

// return token telegram bot or exit
func mustToken() string {
	// Initialized ENV from file
	var err error
	myEnv, err = godotenv.Read()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get the token from the environment variable
	token := myEnv["TELEGRAM_BOT_TOKEN"]
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable not set")
	}

	return token
}

func createBot(token string) *tgbotapi.BotAPI {
	fmt.Println(token)

	// Create a new bot instance
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true // Enable debugging mode

	return bot
}
