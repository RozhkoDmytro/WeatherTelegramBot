package telegram

import (
	"strconv"
	"time"

	"projecttelegrambot/pkg/config"
	"projecttelegrambot/pkg/holiday"
	"projecttelegrambot/pkg/weather"

	"git.foxminded.ua/foxstudent107249/telegrambot"
)

const (
	telegramAPI          = "https://api.telegram.org/bot"
	DefaultHelpStartInfo = `
/start   - get keyboard with flags
/help    - too same like start
/about  - get some information about me
/links   - send my(developer) links`

	DefaultLinksInfo = `
https://www.linkedin.com/in/dmytro-rozhko-bas-1c-golang-junior/
https://animated-panda-0382af.netlify.app/
	`
)

type TelegramService struct {
	config      *config.Config
	apiTelegram *telegrambot.ApiTelegramBot
	apiHoliday  *holiday.ApiHoliday
	apiWeather  *weather.ApiWeather
}

var infoMap = map[string]string{
	"/start": DefaultHelpStartInfo,
	"/help":  DefaultHelpStartInfo,
	"/about": "Rozhko Dmytro; Go developer; bad character; unmarried (C)",
	"/links": DefaultLinksInfo,
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

func NewMyTelegramApp(cfg *config.Config, apiTelegram *telegrambot.ApiTelegramBot, apiHoliday *holiday.ApiHoliday, apiWeather *weather.ApiWeather) *TelegramService {
	return &TelegramService{config: cfg, apiTelegram: apiTelegram, apiHoliday: apiHoliday, apiWeather: apiWeather}
}

func (c *TelegramService) SendResponse(update *telegrambot.Update) error {
	command := update.Message.Text
	chatId := update.Message.Chat.ID

	/* 	if command == "" {
		geotxt := "Latitude: " + strconv.FormatFloat(update.Message.Location.Latitude, 'f', 6, 64) +
			"\nLongitude: " + strconv.FormatFloat(update.Message.Location.Longitude, 'f', 6, 64)
		fmt.Println(geotxt)
	} */
	switch command {
	case "/start":
		_, err := c.apiTelegram.CreateReplyKeyboard(chatId, command, DefualtKeyboard)
		return err

	case "/weather":
		_, err := c.apiTelegram.CreateReplyKeyboard(chatId, "Pls, get location", DefualtKeyboardGeolacation)
		return err
	case "Give Your location":

		geotxt := "Latitude: " + strconv.FormatFloat(update.Message.Location.Latitude, 'f', 6, 64) +
			"\nLongitude: " + strconv.FormatFloat(update.Message.Location.Longitude, 'f', 6, 64)

		/*
			resp, err := c.apiWeather.Load(lat, lon)
			if err != nil {
				return err
			}
			text := c.apiWeather.Description(resp)
			_, err = c.bot.CreateReplayMsg(chatId, text)
			return err */

		_, err := c.apiTelegram.CreateReplayMsg(chatId, geotxt)
		return err
	default:
		if infoMap[command] == "" && flagsCountryMap[command] == "" {
			_, err := c.apiTelegram.CreateReplayMsg(chatId, "")
			return err
		} else if infoMap[command] != "" {
			_, err := c.apiTelegram.CreateReplayMsg(chatId, infoMap[command])
			return err
		} else {
			// send API request and create text message with holidays
			text, err := c.apiHoliday.Names(flagsCountryMap[command], time.Now())
			if err != nil {
				return err
			}
			_, err = c.apiTelegram.CreateReplayMsg(chatId, text)
			return err
		}
	}
}

/* func (c *TelegramService) getGeolocation(chatId int) (float64, float64) {
	// This is a mock function. Replace with actual geolocation fetching logic.
	latitude := 50.4501  // Example latitude (Kyiv)
	longitude := 30.5234 // Example longitude (Kyiv)

	if !c.apiTelegram.Debug {
		c.apiTelegram.CreateReplyKeyboard(chatId, "/weatherGeo", DefualtKeyboardGeolacation)
	}

	return latitude, longitude
} */
