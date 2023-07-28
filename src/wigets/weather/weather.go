package weather

import (
	"fmt"
	"github.com/thiswillgowell/light-controller/src/graphics"
	"github.com/thiswillgowell/light-controller/src/graphics/bitmap/text"
	"github.com/thiswillgowell/light-controller/src/wigets/weather/icons"
	"github.com/thiswillgowell/light-controller/src/wigets/weather/openweather"
	"go.uber.org/zap"
	"golang.org/x/image/colornames"
	"image"
	"image/color"
	"image/draw"
	"math"
	"sync/atomic"
)

var logger = zap.S().With("package", "weather_app")

var weather app
var appState atomic.Int32

const (
	notRunning    int32 = iota
	starting      int32 = iota
	failedToStart int32 = iota
	running       int32 = iota
)

type app struct {
	weatherInfo *openweather.Weather
}

func start() bool {
	if appState.CompareAndSwap(notRunning, starting) {
		weatherInfo, err := openweather.NewWeather()
		if err != nil {
			appState.CompareAndSwap(starting, failedToStart)
			logger.Error("failed to start weather app", zap.Error(err))
			return false
		}
		weather = app{
			weatherInfo: weatherInfo,
		}
		appState.CompareAndSwap(starting, running)
		return true
	}

	return true
}

func DrawCurrentTemp(img draw.Image, fontType text.BitFontType, point image.Point, c color.Color) {
	if !start() {
		return
	}
	currentTemp := weather.weatherInfo.Current().Temp
	text.WriteOnImage(fmt.Sprintf("%.0f°", currentTemp), fontType, c, point, img)
}

func CurrentWeather(img draw.Image, point image.Point) {
	if !start() {
		return
	}
	current := weather.weatherInfo.Current()
	icon := icons.GetWeatherIcon(current.Weather[0].Icon, 36)
	graphics.SaveImage("test-icon", icon)

	graphics.PositionImageAround(point.Add(image.Point{42, 20}), icon, img)
	text.WriteOnImage(fmt.Sprintf("%.0f", current.FeelsLike), text.ExtraLarge, color.White, point.Add(image.Point{4, 6}), img)

	//text.WriteOnImage(fmt.Sprintf("/%.0f", current.Temp), text.Small, color.White, point.Add(image.Point{, 20}), img)
	//writerPos = text.WriteOnImage(fmt.Sprintf("%.0f", current.), text.Medium, colornames.Red, point.Add(image.Point{writerPos.X, 24}), img)
	//writerPos = text.WriteOnImage(fmt.Sprintf("%.0f", current.Temp), text.Medium, colornames.Blue, point.Add(image.Point{writerPos.X, 24}), img)
	graphLength := 48
	Next90MinTemp(img, image.Point{3, point.Y + 30}, 10, graphLength)
	TwoDayHighAndLow(img, image.Point{graphLength + 6, point.Y + 29}, text.ExtraSmall)
	// start pos of last item + height of item + offset
	point.Y += 30 + 10 + 3

	ThreeDayForecast(img, image.Point{point.X + 3, point.Y}, 3)
}

func TwoDayHighAndLow(img draw.Image, pos image.Point, fontType text.BitFontType) {
	forecast := weather.weatherInfo.NextTwoDaysForecast()
	maxTemp := forecast[0].Temp
	minTemp := forecast[0].Temp
	for _, hourlyForecast := range weather.weatherInfo.NextTwoDaysForecast() {
		maxTemp = math.Max(hourlyForecast.Temp, maxTemp)
		minTemp = math.Min(hourlyForecast.Temp, minTemp)
	}
	distance := text.WriteOnImage(fmt.Sprintf("%.0f", maxTemp), fontType, colornames.Red, pos, img)
	pos.Y += distance.Y + 2
	text.WriteOnImage(fmt.Sprintf("%.0f", minTemp), fontType, colornames.Blue, pos, img)

}

func Next90MinTemp(img draw.Image, pos image.Point, height, width int) {
	if !start() {
		return
	}
	forecast := weather.weatherInfo.NextTwoDaysForecast()
	maxTemp := forecast[0].Temp
	minTemp := forecast[0].Temp
	values := []float64{}
	for _, hourlyForecast := range forecast {
		maxTemp = math.Max(hourlyForecast.FeelsLike, maxTemp)
		minTemp = math.Min(hourlyForecast.FeelsLike, minTemp)
		values = append(values, hourlyForecast.FeelsLike)
	}
	yPositions := graphics.MapFloatsToInts(values, minTemp, maxTemp, pos.Y, pos.Y+height, width)
	for i, yPos := range yPositions {
		img.Set(pos.X+i, yPos, color.White)
	}

}

func ThreeDayForecast(img draw.Image, pos image.Point, numDays int) {
	if !start() {
		return
	}
	forecasts := weather.weatherInfo.WeekForecast()

	for i, forecast := range forecasts {
		if i >= numDays {
			break
		}
		iconSize := 17
		icon := icons.GetWeatherIcon(forecast.Weather[0].Icon, iconSize)

		draw.Draw(img, image.Rect(pos.X, pos.Y, pos.X+iconSize, pos.Y+iconSize), icon, image.Point{}, draw.Over)
		text.WriteOnImage(fmt.Sprintf("%.0f %.0f", forecast.FeelsLike.Day, forecast.FeelsLike.Eve), text.ExtraSmall, color.White, pos.Add(image.Point{iconSize + 2, 1}), img)

		pos.Y += iconSize
	}
}
