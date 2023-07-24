package openweather

import (
	_ "embed"
	"encoding/json"
	"fmt"
	owm "github.com/briandowns/openweathermap"
	"go.uber.org/zap"
	"io"
	"os"
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

type UnmarshalItem func([]byte) (interface{}, error)
type UpdateItem func() (interface{}, error)

type jsonCache struct {
	lastUpdated    time.Time
	fileName       string
	unmarshal      UnmarshalItem
	update         UpdateItem
	data           interface{}
	updateInterval time.Duration
}

func NewUpdater(name string, update UpdateItem, unmarshal UnmarshalItem, interval time.Duration) {
	j := jsonCache{
		lastUpdated:    time.Time{},
		fileName:       fmt.Sprintf("%s.json", name),
		unmarshal:      unmarshal,
		update:         update,
		updateInterval: interval,
	}
	j.load()
	j.startUpdate()
}

func (j *jsonCache) load() {
	stats, err := os.Stat(j.fileName)
	if err != nil {
		return
	}
	j.lastUpdated = stats.ModTime()

	f, err := os.Open(j.fileName)
	if err != nil {
		logger.Error("failed to open file", "file_name", j.fileName, zap.Error(err))
		return
	}

	data, err := io.ReadAll(f)
	if err != nil {
		logger.Error("failed to read file", "file_name", j.fileName, zap.Error(err))
		return
	}

	j.data, err = j.unmarshal(data)
	if err != nil {
		logger.Error("failed to unmarshal json file", "file_name", j.fileName, zap.Error(err))
		return
	}
	return
}

func (j *jsonCache) save() {
	data, err := json.MarshalIndent(j.data, "", " ")
	if err != nil {
		logger.Error("failed to marshal data", "file_name", j.fileName, zap.Error(err))
		return
	}
	if err := os.WriteFile(j.fileName, data, 0644); err != nil {
		logger.Error("failed to write file", "file_name", j.fileName, zap.Error(err))
		return
	}
}

func (j *jsonCache) startUpdate() {
	update := func() {
		var err error
		j.data, err = j.update()
		if err != nil {
			logger.Error("failed to update", zap.Error(err))
		}
		j.save()
	}

	// do we need to update right now
	if j.lastUpdated.Add(j.updateInterval).Before(time.Now()) {
		update()
	}
	go func() {
		// should we wait until we start to update
		if time.Now().Before(j.lastUpdated.Add(j.updateInterval)) {
			<-time.After(time.Until(j.lastUpdated.Add(j.updateInterval)))
		}
		ticker := time.NewTicker(j.updateInterval)
		for {
			update()
			<-ticker.C
		}
	}()

}

func NewWeather() (*Weather, error) {
	oneCall, err := owm.NewOneCall("F", "EN", apiKey, []string{owm.ExcludeAlerts})
	if err != nil {
		return nil, err
	}
	weather = &Weather{
		oneCall: oneCall,
	}

	NewUpdater("one_call", func() (interface{}, error) {
		err := weather.oneCall.OneCallByCoordinates(homeCord)
		return weather.oneCall, err
	}, func(bytes []byte) (interface{}, error) {
		data := &owm.OneCallData{}
		err := json.Unmarshal(bytes, data)
		if err != nil {
			return nil, err
		}
		return data, nil
	}, time.Minute*10)

	return weather, nil
}

type forecastWrapper struct {
	Days []owm.Forecast5WeatherList `json:"days"`
}

type Weather struct {
	current  *owm.CurrentWeatherData
	forecast *owm.ForecastWeatherData
	oneCall  *owm.OneCallData
}

func (w *Weather) update() {

}

func (w *Weather) CurrentTemp() owm.Main {
	return weather.current.Main
}

func rainChance() ([]time.Time, float64) {

}
