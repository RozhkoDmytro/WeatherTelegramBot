package bot

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
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

type ApiTelegramBot struct {
	Token  string
	Offset int
	ChatId int
	Logger *slog.Logger
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

func NewBot(t string, l *slog.Logger) *ApiTelegramBot {
	telegramBot := ApiTelegramBot{
		Token:  t,
		Logger: l,
	}
	return &telegramBot
}

func (bot *ApiTelegramBot) CreateResponseToCommand(c string) string {
	result := infoMap[c]
	if result == "" {
		result = "Unknown command!"
	}
	return result
}

func (bot *ApiTelegramBot) GetUpdates() (*GetUpdatesResponse, error) {
	url := fmt.Sprintf("%s%s/getUpdates?offset=%d", telegramAPI, bot.Token, bot.Offset)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return bot.parseTelegramRequest(resp)
}

// parseTelegramRequest handles incoming update from the Telegram web hook
func (bot *ApiTelegramBot) parseTelegramRequest(r *http.Response) (*GetUpdatesResponse, error) {
	var update GetUpdatesResponse
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		bot.Logger.Info("could not decode incoming update", "err", err.Error())
		return nil, err
	}

	return &update, nil
}

// sendTextToTelegramChat sends a text message to the Telegram chat identified by its chat Id
func (bot *ApiTelegramBot) Send(text string) (string, error) {
	var telegramApi string = telegramAPI + bot.Token + "/sendMessage"
	response, err := http.PostForm(
		telegramApi,
		url.Values{
			"chat_id": {strconv.Itoa(bot.ChatId)},
			"text":    {text},
		})
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	bodyBytes, errRead := io.ReadAll(response.Body)
	if errRead != nil {
		return "", err
	}
	bodyString := string(bodyBytes)
	bot.Logger.Info("Body of Telegram Response:", "body", bodyString)

	return bodyString, nil
}
