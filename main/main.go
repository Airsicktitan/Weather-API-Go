package main

import (
	APIkey "druc/sun/apiKey"
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
		Region  string `json:"region"`
		Country string `json:"country"`
	} `json:"location"`

	Current struct {
		TempF     float64 `json:"temp_f"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`

	Forecast struct {
		Forecastday []struct {
			Date string `json:"date"`
			Hour []struct {
				TimeEpoch int64   `json:"time_epoch"`
				TempF     float64 `json:"temp_f"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
				ChanceOfRain float64 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {
	title := "Weather API:"
	fmt.Println(title)
	q := "Pittsburgh,Pennsylvania"

	if len(os.Args) >= 2 {
		q = os.Args[1]
	}

	resp, err := http.Get("http://api.weatherapi.com/v1/forecast.json?key=" + APIkey.APIkey + "&q=" + q + "&days=2&aqi=no&alerts=no")
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	for i := 0; i < len(title); i++ {
		fmt.Print("-")
	}

	fmt.Print("\n")
	fmt.Print("\n")

	if resp.StatusCode != 200 {
		panic("Weather API not available")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		panic(err)
	}

	location, current, hours, date := weather.Location, weather.Current, weather.Forecast.Forecastday[0].Hour, weather.Forecast.Forecastday[0].Date

	fmt.Printf("Current Conditions for the hour:\n\nDate: %s\nLocation: %s, %s\nTemp(F): %.0f°\nCurrent Condition: %s\n\n", date, location.Name, location.Region, current.TempF, current.Condition.Text)
	bottomUnderscore := "Chance of Rain: " + current.Condition.Text

	for i := 0; i < len(bottomUnderscore); i++ {
		fmt.Print("-")
	}

	fmt.Print("\n")
	fmt.Print("\n")

	hourlyForecastHeader := "Hourly forecast:"
	fmt.Println(hourlyForecastHeader)

	for i := 0; i < len(hourlyForecastHeader); i++ {
		fmt.Print("-")
	}

	fmt.Print("\n")

	for _, hour := range hours {
		date := time.Unix(hour.TimeEpoch, 0)
		var message string

		if date.Hour() == time.Now().Hour() {
			message = fmt.Sprintf("%s - %.0f° F, %.0f%%, %s\n", date.Format("15:04"), hour.TempF, hour.ChanceOfRain, hour.Condition.Text)
		} else if date.Before(time.Now()) {
			continue
		}
		message = fmt.Sprintf("%s - %.0f° F, Chance of rain: %.0f%%, %s\n", date.Format("15:04"), hour.TempF, hour.ChanceOfRain, hour.Condition.Text)

		if hour.ChanceOfRain < 40 {
			color.HiGreen(message)
		} else if hour.ChanceOfRain > 60 && hour.ChanceOfRain < 80 {
			color.Yellow(message)
		} else {
			color.HiRed(message)
		}
	}

	fmt.Println("\nPress enter to exit...")
	fmt.Scanln()
}
