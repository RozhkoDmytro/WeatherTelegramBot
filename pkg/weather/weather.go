package weather

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"strings"
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

func NewApiWeather(client *http.Client, url string, t string) *ApiWeather {
	return &ApiWeather{
		token:   t,
		baseURL: url,
		client:  client,
	}
}

func (api *ApiWeather) Load(latitude, longitude float64) (*WeatherResponse, error) {
	url := fmt.Sprintf(api.baseURL+"weather?lat=%s&lon=%s&appid=%s", strconv.FormatFloat(latitude, 'f', 6, 64), strconv.FormatFloat(longitude, 'f', 6, 64), api.token)
	fmt.Println(url)
	resp, err := api.client.Get(url)
	if err != nil {
		return &WeatherResponse{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &WeatherResponse{}, err
	}

	var w WeatherResponse
	if err := json.Unmarshal(body, &w); err != nil {
		return &WeatherResponse{}, err
	}

	return &w, nil
}

func (resp *WeatherResponse) Description() string {
	return formatWeatherResponse(resp)
}

func FormatWeatherResponse2(weatherResponse *WeatherResponse) string {
	const weatherTemplate = `*Current weather* in {{.Name}}:
<b>*Temperature*:</b> {{.Main.Temp}}°C
<b>Pressure:</b> {{.Main.Pressure}} hPa
<b>Humidity:</b> {{.Main.Humidity}}%
<b>Description:</b> {{(index .Weather 0).Description}}`

	tmpl, err := template.New("weather").Parse(weatherTemplate)
	if err != nil {
		return ""
	}

	var buf bytes.Buffer
	writer := io.Writer(&buf)
	err = tmpl.Execute(writer, weatherResponse)
	if err != nil {
		return ""
	}

	return buf.String()
}

// formatWeatherResponse formats the weather data with MarkdownV2
func formatWeatherResponse(weatherResponse *WeatherResponse) string {
	return fmt.Sprintf(
		"*Current weather in %s:*\n"+
			"Temperature: %.2f°C\n"+
			"Pressure: %d hPa\n"+
			"Humidity: %d%%\n"+
			"Description: %s",
		weatherResponse.Name,
		weatherResponse.Main.Temp,
		weatherResponse.Main.Pressure,
		weatherResponse.Main.Humidity,
		escapeMarkdownV2(weatherResponse.Weather[0].Description),
	)
}

func escapeMarkdownV2(text string) string {
	replacer := strings.NewReplacer(
		"_", "\\_", "*", "\\*", "[", "\\[", "]", "\\]", "(",
		"\\(", ")", "\\)", "~", "\\~", "`", "\\`", ">", "\\>",
		"#", "\\#", "+", "\\+", "-", "\\-", "=", "\\=", "|",
		"\\|", "{", "\\{", "}", "\\}", ".", "\\.", "!", "\\!",
	)

	return replacer.Replace(text)
}
