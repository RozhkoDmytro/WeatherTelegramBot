package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestParseUpdateMessageWithText(t *testing.T) {
	chat := Chat{Id: 1}

	msg := Message{
		Text: "hello world",
		Chat: chat,
	}

	update := Update{
		UpdateId: 1,
		Message:  msg,
	}

	requestBody, err := json.Marshal(update)
	fmt.Println(update)
	if err != nil {
		t.Errorf("Failed to marshal update in json, got %s", err.Error())
	}
	req := httptest.NewRequest("POST", "http://myTelegramWebHookHandler.com/secretToken", bytes.NewBuffer(requestBody))

	fmt.Println(req)
	updateToTest, errParse := parseTelegramRequest(req)
	if errParse != nil {
		t.Errorf("Expected a <nil> error, got %s", errParse.Error())
	}
	if *updateToTest != update {
		t.Errorf("Expected update %s, got %s", update.Message.Text, updateToTest.Message.Text)
	}
}

type myFakeService func(*http.Request) (*http.Response, error)

func (s myFakeService) RoundTrip(req *http.Request) (*http.Response, error) {
	return s(req)
}

func TestParseUpdateMessageWithReplaceTransport(t *testing.T) {
	client := &http.Client{
		Transport: myFakeService(func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
				Body: io.NopCloser(strings.NewReader(`{{"update_id":1,"message":{"text":"hello world","chat":{"id":1}}}}`)),
			}, nil
		}),
	}

	req := ????

	updateToTest, errParse := parseTelegramRequest(req)


}
