package holiday

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

const holidayApiUrl = "https://holidays.abstractapi.com/v1/"

type ApiHoliday struct {
	Token  string
	Logger *slog.Logger
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

func NewApiHoliday(t string) *ApiHoliday {
	result := ApiHoliday{
		Token: t,
	}
	return &result
}

func (api *ApiHoliday) Load(country string, day time.Time) ([]Holiday, error) {
	url := fmt.Sprintf(holidayApiUrl+"?api_key=%s&country=%s&year=%d&month=%d&day=%d", api.Token, country, day.Year(), day.Month(), day.Day())

	resp, err := http.Get(url)
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

	return holidays, nil
}
