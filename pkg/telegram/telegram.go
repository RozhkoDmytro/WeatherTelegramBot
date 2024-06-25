package telegram

import (
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

func CreateReplayMsg(chatId int, c string, bot *telegrambot.ApiTelegramBot) ([]byte, error) {
	switch c {
	case "/start":
		return bot.CreateReplyKeyboard(chatId, c, DefualtKeyboard)
	default:
		if infoMap[c] == "" && flagsCountryMap[c] == "" {
			return bot.CreateReplayMsg(chatId, "")
		} else if infoMap[c] != "" {
			return bot.CreateReplayMsg(chatId, infoMap[c])
		} else {
			return bot.CreateReplayMsg(chatId, "Holiday")
			/* 			// send API request and create text message with holidays
			   			apiHoliday := holiday.NewApiHoliday(&http.Client{}, holiday.HolidayApiUrl, cfg.TokenHoliday)
			   			text, err := bot.ApiHolyday.Names(flagsCountryMap[c], time.Now())
			   			if err != nil {
			   				return nil, err
			   			}
			   			return bot.CreateReplayMsg(chatID, text) */
		}
	}
}
