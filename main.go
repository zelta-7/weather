package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
)

type Weather struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
		FeelsLikeC float64 `json:"feelslike_c"`
	} `json:"current"`
	Forecast struct {
		Forecastday []struct {
			MaxtempC  float64 `json:"maxtemp_c"`
			MintempC  float64 `json:"mintemp_c"`
			Condition struct {
				Text string `json:"text"`
			} `json:"condition"`
			Hour []struct {
				TimeEpoch int64   `json:"time_epoch"`
				TempC     float64 `json:"temp_c"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
			} `json:"hour"`
		} `json:"forecastday"`
	}
}

func main() {
	q := "India"

	if len(os.Args) >= 2 {
		q = os.Args[1]
	}

	res, err := http.Get("http://api.weatherapi.com/v1/forecast.json?key=ec180872243c4f57a4f153631230105&q&q=" + q + "&days=1&aqi=no&alerts=no")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		panic("API request failed")
	}

	apiResult, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var weather Weather
	err = json.Unmarshal(apiResult, &weather)
	if err != nil {
		panic(err)
	}

	location, current, forecast := weather.Location, weather.Current, weather.Forecast.Forecastday[0].Hour

	fmt.Printf("%s ", location.Country)
	fmt.Printf("%.0f°C, %s\n", current.TempC, current.Condition.Text)
	fmt.Printf("Feels Like: %.2f°C\n", current.FeelsLikeC)
	// fmt.Printf("Max: %.0f°C, Min: %.2f°C, %s\n", weather.Forecast.Forecastday[0].MaxtempC, weather.Forecast.Forecastday[0].MintempC, weather.Forecast.Forecastday[0].Condition.Text)
	fmt.Printf("Hourly Forecast:\n")

	for _, hour := range forecast {
		date := time.Unix(hour.TimeEpoch, 0)

		if date.Before(time.Now()) {
			continue
		}

		message := fmt.Sprintf("%s: %.0f°C, %s\n", date.Format("15:04"), hour.TempC, hour.Condition.Text)

		if hour.TempC < 10 {
			color.Cyan(message)
		} else if hour.TempC > 25 {
			color.Yellow(message)
		} else if hour.TempC > 35 {
			color.Red(message)
		}
	}
}
