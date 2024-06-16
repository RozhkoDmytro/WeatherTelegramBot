package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/joho/godotenv"
)

// Pass token and sensible APIs through environment variables
const (
	telegramApiBaseUrl     string = "https://api.telegram.org/bot"
	telegramApiSendMessage string = "/sendMessage"
	telegramTokenEnv       string = "TELEGRAM_BOT_TOKEN"
)

// Update is a Telegram object that the handler receives every time an user interacts with the bot.
type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

// Message is a Telegram object that can be found in an update.
type Message struct {
	Text string `json:"text"`
	Chat Chat   `json:"chat"`
}

// A Telegram Chat indicates the conversation to which the message belongs.
type Chat struct {
	Id int `json:"id"`
}

var (
	myEnv   map[string]string
	infoMap map[string]string
	token   string
)

// Set all parametrs
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

	// Set token
	token = mustToken()
}

func main() {
}

// parseTelegramRequest handles incoming update from the Telegram web hook
func parseTelegramRequest(r *http.Request) (*Update, error) {
	var update Update
	fmt.Println(r.Body)
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		log.Printf("could not decode incoming update %s", err.Error())
		return nil, err
	}
	return &update, nil
}

// HandleTelegramWebHook sends a message back to the chat with a punchline starting by the message provided by the user.
func HandleTelegramWebHook(w http.ResponseWriter, r *http.Request) {
	// Parse incoming request
	update, err := parseTelegramRequest(r)
	if err != nil {
		log.Printf("error parsing update, %s", err.Error())
		return
	}

	// If we got a message
	if update.Message.Text != "" {
		sendTextToTelegramChat(update.Message.Chat.Id, infoMap[update.Message.Text])
	}

	// Send the punchline back to Telegram
	msg := infoMap[update.Message.Text]
	telegramResponseBody, errTelegram := sendTextToTelegramChat(update.Message.Chat.Id, msg)
	if errTelegram != nil {
		log.Printf("got error %s from telegram, reponse body is %s", errTelegram.Error(), telegramResponseBody)
	} else {
		log.Printf("punchline %s successfuly distributed to chat id %d", msg, update.Message.Chat.Id)
	}
}

// sendTextToTelegramChat sends a text message to the Telegram chat identified by its chat Id
func sendTextToTelegramChat(chatId int, text string) (string, error) {
	log.Printf("Sending %s to chat_id: %d", text, chatId)
	var telegramApi string = apiUrl() + "/sendMessage"
	response, err := http.PostForm(
		telegramApi,
		url.Values{
			"chat_id": {strconv.Itoa(chatId)},
			"text":    {text},
		})
	if err != nil {
		log.Printf("error when posting text to the chat: %s", err.Error())
		return "", err
	}
	defer response.Body.Close()

	bodyBytes, errRead := io.ReadAll(response.Body)
	if errRead != nil {
		log.Printf("error in parsing telegram answer %s", errRead.Error())
		return "", err
	}
	bodyString := string(bodyBytes)
	log.Printf("Body of Telegram Response: %s", bodyString)

	return bodyString, nil
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

func apiUrl() string {
	return telegramApiBaseUrl + token
}
