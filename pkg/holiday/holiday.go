package holiday

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const HolidayApiUrl = "https://holidays.abstractapi.com/v1/"

type ApiHoliday struct {
	Client  *http.Client
	baseURL string
	Token   string
}

// Define a struct to hold the holiday data
type Holiday struct {
	Name        string `json:"name"`
	NameLocal   string `json:"name_local"`
	Language    string `json:"language"`
	Description string `json:"description"`
	Country     string `json:"country"`
	Location    string `json:"location"`
	Type        string `json:"type"`
	Date        string `json:"date"`
	DateYear    string `json:"date_year"`
	DateMonth   string `json:"date_month"`
	DateDay     string `json:"date_day"`
	WeekDay     string `json:"week_day"`
}

func NewApiHoliday(client *http.Client, url string, t string) *ApiHoliday {
	result := ApiHoliday{
		Token:   t,
		baseURL: url,
		Client:  client,
	}
	return &result
}

func (api *ApiHoliday) Load(country string, day time.Time) ([]Holiday, error) {
	url := fmt.Sprintf(api.baseURL+"?api_key=%s&country=%s&year=%d&month=%d&day=%d", api.Token, country, day.Year(), day.Month(), day.Day())

	resp, err := api.Client.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(body))
	var holidays []Holiday
	if err := json.Unmarshal(body, &holidays); err != nil {
		return nil, err
	}

	fmt.Println(holidays)

	return holidays, nil
}

func (api *ApiHoliday) Names(country string, day time.Time) (string, error) {
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
