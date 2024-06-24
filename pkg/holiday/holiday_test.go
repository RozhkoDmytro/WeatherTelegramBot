package holiday

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"projecttelegrambot/pkg/config"

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

func TestNames(t *testing.T) {
	testCases = map[string]GenericTestCase{
		"test 06/24/2024": {
			input:    `[{"name":"Orthodox Pentecost holiday","name_local":"","language":"","description":"","country":"UA","location":"Ukraine","type":"National","date":"06/24/2024","date_year":"2024","date_month":"06","date_day":"24","week_day":"Monday"}]`,
			expected: `[{Orthodox Pentecost holiday    UA Ukraine National 06/24/2024 2024 06 24 Monday}]`,
		},
		"no holidays": {
			input:    `[]`,
			expected: `[]`,
		},
	}

	// Get config with env
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	for name, tc := range testCases {
		// Start a local HTTP server
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			// Send response to be tested
			rw.Write([]byte(tc.input))
		}))

		// Close the server when test finishes
		defer server.Close()

		// Use Client & URL from our local test server
		api := NewApiHoliday(server.Client(), server.URL, cfg.TokenHoliday)
		holidays, err := api.Load("UA", time.Date(2024, 6, 24, 0, 0, 0, 0, time.UTC))
		fmt.Println(err)
		assert.Equal(t, tc.expected, fmt.Sprintf("%s", holidays), name)
	}

	// Test error
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		// assert.Equal(t, req.URL.String(), HolidayApiUrl)
		// Send response to be tested
		rw.Write([]byte(`OK`))
	}))

	// Close the server when test finishes
	defer server.Close()

	// Use Client & URL from our local test server
	api := NewApiHoliday(server.Client(), server.URL, cfg.TokenHoliday)
	_, err = api.Load("UA", time.Date(2024, 6, 24, 0, 0, 0, 0, time.UTC))

	assert.Equal(t, true, err != nil, "Test error Unmarshaling")
}
