package telegram

import (
	"log/slog"
	"time"

	"projecttelegrambot/pkg/config"
	"projecttelegrambot/pkg/holiday"

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

type MyApp struct {
	bot        *telegrambot.ApiTelegramBot
	apiHoliday *holiday.ApiHoliday
	config     *config.Config
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

func NewMyTelegramApp(cfg *config.Config, bot *telegrambot.ApiTelegramBot, apiHoliday *holiday.ApiHoliday) (*MyApp, error) {
	return &MyApp{bot: bot, apiHoliday: apiHoliday, config: cfg}, nil
}

func (c *MyApp) SendResponse(update *telegrambot.Update) error {
	command := update.Message.Text
	chatId := update.Message.Chat.ID

	switch command {
	case "/start":
		_, err := c.bot.CreateReplyKeyboard(chatId, command, DefualtKeyboard)
		return err
	default:
		if infoMap[command] == "" && flagsCountryMap[command] == "" {
			_, err := c.bot.CreateReplayMsg(chatId, "")
			return err
		} else if infoMap[command] != "" {
			_, err := c.bot.CreateReplayMsg(chatId, infoMap[command])
			return err
		} else {
			// send API request and create text message with holidays
			text, err := c.apiHoliday.Names(flagsCountryMap[command], time.Now())
			if err != nil {
				return err
			}
			_, err = c.bot.CreateReplayMsg(chatId, text)
			return err
		}
	}
}

func (c *MyApp) SetLogger(l *slog.Logger) {
	c.bot.Logger = l
}
