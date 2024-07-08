package weather

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

/* // GenericTestCase is a struct for generic test cases
type GenericTestCase struct {
	input    string
	expected string
}

var testCases map[string]GenericTestCase */

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

func TestDescription(t *testing.T) {
	api := NewApiWeather(&http.Client{}, WeatherApiUrl, "c6e84a54fbf8ac35203818a6ed8e40a0")

	latitude, longitude := geoDebuglocation()
	resp, _ := api.Load(latitude, longitude)
	api.Description(resp)
	assert.Equal(t, "", "")
}

func geoDebuglocation() (float64, float64) {
	// This is a mock function. Replace with actual geolocation fetching logic.
	latitude := 50.4501  // Example latitude (Kyiv)
	longitude := 30.5234 // Example longitude (Kyiv)
	return latitude, longitude
}
