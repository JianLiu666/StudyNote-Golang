package main

import (
	"context"
	"fmt"

	"github.com/hectormalot/omgo"
)

func main() {
	c, _ := omgo.NewClient()

	loc, _ := omgo.NewLocation(25.0648093, 121.5320234)

	current, _ := c.CurrentWeather(context.Background(), loc, &omgo.Options{
		Timezone: "Asia/Taipei",
	})
	fmt.Println("The temperature in 台北市中山區 is:", current.Temperature)

	forecast, _ := c.Forecast(context.Background(), loc, &omgo.Options{
		Timezone:      "Asia/Taipei",
		PastDays:      7,
		HourlyMetrics: []string{"cloudcover, relativehumidity_2m"},
		DailyMetrics:  []string{"temperature_2m_max"},
	})
	fmt.Println(forecast)
}
