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
)

const (
	telegramAPI = "https://api.telegram.org/bot"
)

func RunCommand(update types.Update) string {
	fmt.Println("212341234134 " + update.Message.Text)

	helpStartInfo := `
/start   - get information about all bot commands
/help    - too same like start
/about  - get some information about me
/links   - send my(developer) links`

	linksInfo := `
https://www.linkedin.com/in/dmytro-rozhko-bas-1c-golang-junior/
https://animated-panda-0382af.netlify.app/
	`
	infoMap := map[string]string{
		"/start": helpStartInfo,
		"/help":  helpStartInfo,
		"/about": "Rozhko Dmytro; Go developer; bad character; unmarried (C)",
		"/links": linksInfo,
	}

	return infoMap[update.Message.Text]
}

func GetUpdates(token string, offset int) (*types.GetUpdatesResponse, error) {
	url := fmt.Sprintf("%s%s/getUpdates?offset=%d", telegramAPI, token, offset)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var updates types.GetUpdatesResponse
	if err := json.NewDecoder(resp.Body).Decode(&updates); err != nil {
		return nil, err
	}
	return &updates, nil
}

// sendTextToTelegramChat sends a text message to the Telegram chat identified by its chat Id
func Send(chatId int, text string, token string) (string, error) {
	log.Printf("Sending %s to chat_id: %d", text, chatId)
	var telegramApi string = telegramAPI + token + "/sendMessage"
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
