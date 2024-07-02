package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const WeatherApiUrl = "https://api.openweathermap.org/data/2.5/"

type ApiWeather struct {
	client  *http.Client
	baseURL string
	token   string
}

// Structs to match the JSON response

type WeatherResponse struct {
	Coord      Coord     `json:"coord"`
	Weather    []Weather `json:"weather"`
	Base       string    `json:"base"`
	Main       Main      `json:"main"`
	Visibility int       `json:"visibility"`
	Wind       Wind      `json:"wind"`
	Clouds     Clouds    `json:"clouds"`
	Dt         int64     `json:"dt"`
	Sys        Sys       `json:"sys"`
	Timezone   int       `json:"timezone"`
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Cod        int       `json:"cod"`
}

type Coord struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

type Weather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type Main struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	Pressure  int     `json:"pressure"`
	Humidity  int     `json:"humidity"`
}

type Wind struct {
	Speed float64 `json:"speed"`
	Deg   float64 `json:"deg"`
}

type Clouds struct {
	All int `json:"all"`
}

type Sys struct {
	Type    int     `json:"type"`
	ID      int     `json:"id"`
	Message float64 `json:"message"`
	Country string  `json:"country"`
	Sunrise int64   `json:"sunrise"`
	Sunset  int64   `json:"sunset"`
}

func NewApiHoliday(client *http.Client, url string, t string) *ApiWeather {
	result := ApiWeather{
		token:   t,
		baseURL: url,
		client:  client,
	}
	return &result
}

func (api *ApiWeather) Load(country string, day time.Time) ([]WeatherResponse, error) {
	url := fmt.Sprintf(api.baseURL+"weather?lat={%s}&lon={%s}&appid={%s}", country, day.Year(), day.Month(), day.Day(), api.token)

	resp, err := api.client.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var holidays []WeatherResponse
	if err := json.Unmarshal(body, &holidays); err != nil {
		return nil, err
	}

	return holidays, nil
}

func (api *ApiWeather) CreateWeatherDescription(country string, day time.Time) (string, error) {
	holydays, err := api.Load(country, day)
	if err != nil {
		return "", err
	}
	text := ""
	for _, h := range holydays {
		text += h.Name
	}

	if text == "" {
		text = "There are no holidays in this country today, so sad !"
	}
	return text, nil
}
