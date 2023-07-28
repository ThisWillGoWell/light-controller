package openweather

import (
	_ "embed"
	"encoding/json"
	owm "github.com/briandowns/openweathermap"
	"github.com/thiswillgowell/light-controller/updater"
	"go.uber.org/zap"
	"time"
)

var logger = zap.S().With("package", "weather")

//go:embed key.txt
var apiKey string

var weather *Weather
var homeCord = &owm.Coordinates{
	Latitude:  41.94970638786151,
	Longitude: -87.66443131375526,
}

func NewWeather() (*Weather, error) {
	oneCall, err := owm.NewOneCall("F", "EN", apiKey, []string{owm.ExcludeAlerts})
	if err != nil {
		return nil, err
	}
	weather = &Weather{
		oneCall: oneCall,
	}

	weather.updater = updater.NewUpdater("one_call", func() (interface{}, error) {
		err := weather.oneCall.OneCallByCoordinates(homeCord)
		return weather.oneCall, err
	}, func(bytes []byte) (interface{}, error) {
		err := json.Unmarshal(bytes, weather.oneCall)
		if err != nil {
			return nil, err
		}
		return weather.oneCall, nil
	}, time.Minute*10)

	return weather, nil
}

type forecastWrapper struct {
	Days []owm.Forecast5WeatherList `json:"days"`
}

type Weather struct {
	oneCall *owm.OneCallData
	updater *updater.ApiJsonUpdater
}

func (w *Weather) Stop() {
	w.updater.Stop()
}

func (w *Weather) Current() owm.OneCallCurrentData {
	return weather.oneCall.Current
}

type TwoDayHourlyForecast struct {
	Temps   []owm.OneCallHourlyData
	MaxTemp float64
	MinTemp float64
}

func (w *Weather) NextTwoDaysForecast() []owm.OneCallHourlyData {
	return w.oneCall.Hourly
}

func (w *Weather) WeekForecast() []owm.OneCallDailyData {
	// ignore current day
	return w.oneCall.Daily
}
