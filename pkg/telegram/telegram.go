package telegram

import (
	"time"

	"projecttelegrambot/pkg/holiday"
	"projecttelegrambot/pkg/weather"

	"git.foxminded.ua/foxstudent107249/telegrambot"
)

const (
	telegramAPI          = "https://api.telegram.org/bot"
	DefaultHelpStartInfo = `
/start   - get keyboard with flags
/help    - too same like start
/weather - get the current weather for your location
/about  - get some information about me
/links   - send my(developer) links`

	DefaultLinksInfo = `
https://www.linkedin.com/in/dmytro-rozhko-bas-1c-golang-junior/
https://animated-panda-0382af.netlify.app/
	`
	CurrentParseMode = telegrambot.ModeMarkdownV2
)

type TelegramService struct {
	apiTelegram *telegrambot.ApiTelegramBot
	apiHoliday  *holiday.ApiHoliday
	apiWeather  *weather.ApiWeather
}

var infoMap = map[string]string{
	"/start":   DefaultHelpStartInfo,
	"/help":    DefaultHelpStartInfo,
	"/weather": "Get the current weather for your location",
	"/about":   "Rozhko Dmytro; Go developer; bad character; unmarried (C)",
	"/links":   DefaultLinksInfo,
}

var DefualtKeyboard = telegrambot.ReplyKeyboardMarkup{
	Keyboard: [][]telegrambot.KeyboardButton{
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

var DefualtKeyboardGeolacation = telegrambot.ReplyKeyboardMarkup{
	Keyboard: [][]telegrambot.KeyboardButton{
		{
			{Text: "Give Your location", RequestLocation: true},
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

func NewMyTelegramService(apiTelegram *telegrambot.ApiTelegramBot, apiHoliday *holiday.ApiHoliday, apiWeather *weather.ApiWeather) *TelegramService {
	return &TelegramService{apiTelegram: apiTelegram, apiHoliday: apiHoliday, apiWeather: apiWeather}
}

func (c *TelegramService) CreateSendResponse(update *telegrambot.Update) error {
	command := update.Message.Text
	chatId := update.Message.Chat.ID

	switch command {
	case "/start":
		_, err := c.apiTelegram.CreateReplyKeyboard(chatId, command, DefualtKeyboard)
		return err
	case "/weather":
		_, err := c.apiTelegram.CreateReplyKeyboard(chatId, "Pls, get location", DefualtKeyboardGeolacation)
		return err
	default:
		// Unknown
		if isUnknownCommand(update) {
			_, err := c.apiTelegram.CreateReplayMsg(chatId, "", CurrentParseMode)
			return err
		}
		// standart command
		if infoMap[command] != "" {
			_, err := c.apiTelegram.CreateReplayMsg(chatId, infoMap[command], CurrentParseMode)
			return err
		}
		// responce with country value
		if flagsCountryMap[command] != "" {
			return c.createReplayMsgHoliday(update)
		}
		// responce with fill Location
		if update.Message.Location != nil {
			return c.createReplayMsgWeather(update)
		}
	}
	return nil
}

// send API request and create text message with weather and sen him
func (c *TelegramService) createReplayMsgWeather(update *telegrambot.Update) error {
	chatId := update.Message.Chat.ID
	resp, err := c.apiWeather.Load(update.Message.Location.Latitude, update.Message.Location.Longitude)
	if err != nil {
		return err
	}
	geotxt, err := resp.Description()
	if err != nil {
		return err
	}
	_, err = c.apiTelegram.CreateReplayMsg(chatId, geotxt, CurrentParseMode)
	return err
}

// send API request and create text message with holidays and sen him
func (c *TelegramService) createReplayMsgHoliday(update *telegrambot.Update) error {
	command := update.Message.Text
	chatId := update.Message.Chat.ID
	text, err := c.apiHoliday.Names(flagsCountryMap[command], time.Now())
	if err != nil {
		return err
	}
	_, err = c.apiTelegram.CreateReplayMsg(chatId, text, CurrentParseMode)
	return err
}

func isUnknownCommand(update *telegrambot.Update) bool {
	command := update.Message.Text
	return infoMap[command] == "" && flagsCountryMap[command] == "" && update.Message.Location == nil
}
