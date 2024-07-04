package telegram

import (
	"io"
	"log/slog"
	"net/http"
	"os"
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
	Bot        *telegrambot.ApiTelegramBot
	ApiHoliday *holiday.ApiHoliday
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

func NewMyApp() (*MyApp, error) {
	// Get config with env
	cfg, err := config.Load()
	if err != nil {
		return &MyApp{}, err
	}

	// Create logger
	logger, err := createLogger(cfg.NameLog)
	if err != nil {
		return &MyApp{}, err
	}
	// Create a new telegram bot
	bot, err := telegrambot.NewBot(cfg.Token, logger)
	if err != nil {
		return &MyApp{}, err
	}

	apiHoliday := holiday.NewApiHoliday(&http.Client{}, holiday.HolidayApiUrl, cfg.TokenHoliday)

	return &MyApp{Bot: bot, ApiHoliday: apiHoliday, config: &cfg}, nil
}

// Create logger and set fields
func createLogger(NameLog string) (*slog.Logger, error) {
	// Create logger
	file, err := os.OpenFile(NameLog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		return nil, err
	}

	w := io.MultiWriter(os.Stderr, file)
	handler := slog.NewJSONHandler(w, &slog.HandlerOptions{
		AddSource: true,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger, nil
}

func (c *MyApp) SendResponse(update *telegrambot.Update) error {
	command := update.Message.Text
	chatId := update.Message.Chat.ID

	switch command {
	case "/start":
		return c.Bot.CreateReplyKeyboard(chatId, command, DefualtKeyboard)
	default:
		if infoMap[command] == "" && flagsCountryMap[command] == "" {
			return c.Bot.CreateReplayMsg(chatId, "")
		} else if infoMap[command] != "" {
			return c.Bot.CreateReplayMsg(chatId, infoMap[command])
		} else {
			// send API request and create text message with holidays
			text, err := c.ApiHoliday.Names(flagsCountryMap[command], time.Now())
			if err != nil {
				return err
			}
			return c.Bot.CreateReplayMsg(chatId, text)
		}
	}
}
