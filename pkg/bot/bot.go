package bot

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type Update struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		MessageID int `json:"message_id"`
		From      struct {
			ID        int    `json:"id"`
			IsBot     bool   `json:"is_bot"`
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
			Language  string `json:"language_code"`
		} `json:"from"`
		Chat struct {
			ID        int64  `json:"id"`
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
			Type      string `json:"type"`
		} `json:"chat"`
		Date     int    `json:"date"`
		Text     string `json:"text"`
		Entities []struct {
			Offset int    `json:"offset"`
			Length int    `json:"length"`
			Type   string `json:"type"`
		} `json:"entities"`
	} `json:"message"`
}

type GetUpdatesResponse struct {
	OK     bool     `json:"ok"`
	Result []Update `json:"result"`
}

const (
	telegramAPI          = "https://api.telegram.org/bot"
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

func GetUpdates(token string, offset int) (*GetUpdatesResponse, error) {
	url := fmt.Sprintf("%s%s/getUpdates?offset=%d", telegramAPI, token, offset)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ParseTelegramRequest(resp)
}

// parseTelegramRequest handles incoming update from the Telegram web hook
func ParseTelegramRequest(r *http.Response) (*GetUpdatesResponse, error) {
	var update GetUpdatesResponse
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
