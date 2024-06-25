package telegram

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"git.foxminded.ua/foxstudent107249/telegrambot"
	"github.com/stretchr/testify/assert"
)

// GenericTestCase is a struct for generic test cases
type GenericTestCase struct {
	input    string
	expected string
}

var testCases map[string]GenericTestCase

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestParseTelegramRequest(t *testing.T) {
	expectString := `{"ok":false,"result":[{"update_id":0,"message":{"message_id":0,"from":{"id":0,"is_bot":false,"first_name":"","username":"","language_code":""},"chat":{"id":0,"first_name":"","username":"","type":""},"date":0,"text":"/start","entities":null}}]}`

	client := NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
		return &http.Response{
			StatusCode: 200,
			// Send response to be tested
			Body: io.NopCloser(strings.NewReader(expectString)),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})

	resp, _ := client.Get("")
	var telegramBot telegrambot.ApiTelegramBot
	result, _ := telegramBot.ParseTelegramRequest(resp)
	assert.Equal(t, "/start", result.Result[0].Message.Text)
}

func TestCreateReplayMsg(t *testing.T) {
	var telegramBot telegrambot.ApiTelegramBot
	chatId := 123456

	testCases = map[string]GenericTestCase{
		"test command start": {
			input:    "/start",
			expected: `{"chat_id":123456,"text":"/start","reply_markup":{"keyboard":[[{"text":"🇺🇸 USA"},{"text":"🇬🇧 UK"},{"text":"🇨🇦 Canada"}],[{"text":"🇦🇺 Australia"},{"text":"🇮🇳 India"},{"text":"🇺🇦 Ukraine"}]],"resize_keyboard":true,"one_time_keyboard":true}}`,
		},
		"test command info": {
			input:    "/help",
			expected: DefaultHelpStartInfo,
		},
		"test command about": {
			input:    "/about",
			expected: "Rozhko Dmytro; Go developer; bad character; unmarried (C)",
		},
		"test command links": {
			input:    "/links",
			expected: DefaultLinksInfo,
		},
	}

	for name, tc := range testCases {
		t.Run("Test Create Replay Msg", func(t *testing.T) {
			// Keep

			if tc.input == "/start" {
				expected := []byte(tc.expected)
				result, _ := telegramBot.CreateReplyKeyboard(chatId, tc.input, DefualtKeyboard)
				assert.Equal(t, expected, result, name)
			}
		})
	}
}
