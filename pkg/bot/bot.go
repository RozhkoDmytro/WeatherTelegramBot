package bot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"projecttelegrambot/pkg/holiday"
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
			ID        int    `json:"id"`
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
	Token      string
	Offset     int
	Logger     *slog.Logger
	ApiHolyday *holiday.ApiHoliday
}

type ReplyKeyboardMarkup struct {
	Keyboard        [][]KeyboardButton `json:"keyboard"`
	ResizeKeyboard  bool               `json:"resize_keyboard"`
	OneTimeKeyboard bool               `json:"one_time_keyboard"`
}

type KeyboardButton struct {
	Text string `json:"text"`
}

type ReplayMsg struct {
	ChatId int    `json:"chat_id"`
	Text   string `json:"text"`
}

type ReplayKeyboardMsg struct {
	ChatID      int                 `json:"chat_id"`
	Text        string              `json:"text"`
	ReplyMarkup ReplyKeyboardMarkup `json:"reply_markup"`
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

var DefualtKeyboard = ReplyKeyboardMarkup{
	Keyboard: [][]KeyboardButton{
		{
			{Text: DefaultFlags[0]},
			{Text: DefaultFlags[1]},
			{Text: DefaultFlags[2]},
		},
		{
			{Text: DefaultFlags[3]},
			{Text: DefaultFlags[4]},
			{Text: DefaultFlags[5]},
		},
	},
	ResizeKeyboard:  true,
	OneTimeKeyboard: true,
}

var DefaultFlags = []string{
	"🇺🇸 USA",
	"🇬🇧 UK",
	"🇨🇦 Canada",
	"🇦🇺 Australia",
	"🇮🇳 India",
	"🇺🇦 Ukraine",
}

var flagsCountryMap = map[string]string{
	DefaultFlags[0]: "US",
	DefaultFlags[1]: "GB",
	DefaultFlags[2]: "CA",
	DefaultFlags[3]: "AU",
	DefaultFlags[4]: "IN",
	DefaultFlags[5]: "UA",
}

func NewBot(t string, l *slog.Logger, a *holiday.ApiHoliday) *ApiTelegramBot {
	telegramBot := ApiTelegramBot{
		Token:      t,
		Logger:     l,
		ApiHolyday: a,
	}
	return &telegramBot
}

func (bot *ApiTelegramBot) CreateResponseToCommand(chatId int, c string) ([]byte, error) {
	switch c {
	case "/start":
		return bot.createReplyKeyboard(chatId, c)
	default:
		return bot.createReplayMsg(chatId, c)
	}
}

func (bot *ApiTelegramBot) GetUpdates() (*GetUpdatesResponse, error) {
	url := fmt.Sprintf("%s%s/getUpdates?offset=%d", telegramAPI, bot.Token, bot.Offset)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		err := errors.New("Error: received status code:" + strconv.Itoa(resp.StatusCode))
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
func (bot *ApiTelegramBot) Send(chatId int, body []byte) (string, error) {
	var urlTelegram string = telegramAPI + bot.Token + "/sendMessage"
	resp, err := http.Post(urlTelegram, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := errors.New("Error: received status code:" + strconv.Itoa(resp.StatusCode) + " " + string(body))
		return "", err
	}

	bodyBytes, errRead := io.ReadAll(resp.Body)
	if errRead != nil {
		return "", err
	}
	bodyString := string(bodyBytes)
	bot.Logger.Info("Body of Telegram Response:", "body", bodyString)

	return bodyString, nil
}

// Create Keyboard with defualt key
func (bot *ApiTelegramBot) createReplyKeyboard(chatID int, text string) ([]byte, error) {
	msg := ReplayKeyboardMsg{
		ChatID:      chatID,
		Text:        text,
		ReplyMarkup: DefualtKeyboard,
	}

	body, err := json.Marshal(msg)
	fmt.Println(string(body))
	if err != nil {
		return nil, err
	}

	return body, nil
}

// when we know the command shudn't create Keyboard we try send text message
func (bot *ApiTelegramBot) createReplayMsg(chatID int, c string) ([]byte, error) {
	if infoMap[c] == "" && flagsCountryMap[c] == "" {
		return bot.createMsgBody(chatID, "")
	} else if infoMap[c] != "" {
		return bot.createMsgBody(chatID, infoMap[c])
	} else {
		// send API request and create text message with holidays
		text, err := bot.ApiHolyday.Names(flagsCountryMap[c], time.Now())
		if err != nil {
			return nil, err
		}
		return bot.createMsgBody(chatID, text)
	}
}

// helpful func to transformation text in json []byte
func (bot *ApiTelegramBot) createMsgBody(chatID int, text string) ([]byte, error) {
	if text == "" {
		text = "Unknown command!"
	}

	msg := ReplayMsg{
		ChatId: chatID,
		Text:   text,
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return body, nil
}
