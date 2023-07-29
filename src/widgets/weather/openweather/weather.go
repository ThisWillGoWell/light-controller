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

var weather *openWeather
var homeCord = &owm.Coordinates{
	Latitude:  41.94970638786151,
	Longitude: -87.66443131375526,
}

func init() {
	var err error
	weather, err = newWeather()
	if err != nil {
		panic(err)
	}
}

type openWeather struct {
	oneCall *owm.OneCallData
	updater *updater.ApiJsonUpdater
}

func newWeather() (*openWeather, error) {
	oneCall, err := owm.NewOneCall("F", "EN", apiKey, []string{owm.ExcludeAlerts})
	if err != nil {
		return nil, err
	}
	weather = &openWeather{
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

func Current() owm.OneCallCurrentData {
	return weather.oneCall.Current
}

type TwoDayHourlyForecast struct {
	Temps   []owm.OneCallHourlyData
	MaxTemp float64
	MinTemp float64
}

func NextTwoDaysForecast() []owm.OneCallHourlyData {
	return weather.oneCall.Hourly
}

func WeekForecast() []owm.OneCallDailyData {
	return weather.oneCall.Daily
}
