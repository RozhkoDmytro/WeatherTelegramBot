package bot

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"projecttelegrambot/pkg/types"

	"github.com/joho/godotenv"
)

const (
	telegramAPI          = "https://api.telegram.org/bot"
	telegramTokenEnv     = "TELEGRAM_BOT_TOKEN"
	DefaultHelpStartInfo = `
/start   - get information about all bot commands
/help    - too same like start
/about  - get some information about me
/links   - send my(developer) links`

	DefaultLinksInfo = `
https://www.linkedin.com/in/dmytro-rozhko-bas-1c-golang-junior/
https://animated-panda-0382af.netlify.app/
	`
)

var infoMap = map[string]string{
	"/start": DefaultHelpStartInfo,
	"/help":  DefaultHelpStartInfo,
	"/about": "Rozhko Dmytro; Go developer; bad character; unmarried (C)",
	"/links": DefaultLinksInfo,
}

func GetInfo(c string) string {
	result := infoMap[c]
	if result == "" {
		result = "Unknown command!"
	}
	return result
}

func GetUpdates(token string, offset int) (*types.GetUpdatesResponse, error) {
	url := fmt.Sprintf("%s%s/getUpdates?offset=%d", telegramAPI, token, offset)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ParseTelegramRequest(resp)
}

// parseTelegramRequest handles incoming update from the Telegram web hook
func ParseTelegramRequest(r *http.Response) (*types.GetUpdatesResponse, error) {
	var update types.GetUpdatesResponse
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		log.Printf("could not decode incoming update %s", err.Error())
		return nil, err
	}

	return &update, nil
}

// sendTextToTelegramChat sends a text message to the Telegram chat identified by its chat Id
func Send(chatId int, text string, telegramToken string) (string, error) {
	log.Printf("Sending %s to chat_id: %d", text, chatId)
	var telegramApi string = telegramAPI + telegramToken + "/sendMessage"
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
func MustToken() string {
	// Initialized ENV from file
	myEnv, err := godotenv.Read()
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
