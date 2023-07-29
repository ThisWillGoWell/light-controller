package weather

import (
	"fmt"
	"github.com/thiswillgowell/light-controller/src/graphics"
	"github.com/thiswillgowell/light-controller/src/graphics/bitmap/text"
	"github.com/thiswillgowell/light-controller/src/widgets/weather/icons"
	"github.com/thiswillgowell/light-controller/src/widgets/weather/openweather"
	"go.uber.org/zap"
	"golang.org/x/image/colornames"
	"image"
	"image/color"
	"image/draw"
	"math"
)

var logger = zap.S().With("package", "weather_app")

func DrawCurrentTemp(img draw.Image, fontType text.BitFontType, point image.Point, c color.Color) {
	currentTemp := openweather.Current().Temp
	text.WriteOnImage(fmt.Sprintf("%.0f°", currentTemp), fontType, c, point, img)
}

func CurrentTemp(img draw.Image, point image.Point) {
	current := openweather.Current()
	icon := icons.GetWeatherIcon(current.Weather[0].Icon, 36)
	graphics.PositionImageAround(point.Add(image.Point{X: 42, Y: 20}), icon, img)
	text.WriteOnImage(fmt.Sprintf("%.0f", current.FeelsLike), text.ExtraLarge, color.White, point.Add(image.Point{X: 4, Y: 6}), img)
}

func TwoDayHighAndLow(img draw.Image, pos image.Point, fontType text.BitFontType) {
	forecast := openweather.NextTwoDaysForecast()
	maxTemp := forecast[0].Temp
	minTemp := forecast[0].Temp
	for _, hourlyForecast := range openweather.NextTwoDaysForecast() {
		maxTemp = math.Max(hourlyForecast.Temp, maxTemp)
		minTemp = math.Min(hourlyForecast.Temp, minTemp)
	}
	distance := text.WriteOnImage(fmt.Sprintf("%.0f", maxTemp), fontType, colornames.Red, pos, img)
	pos.Y += distance.Y + 2
	text.WriteOnImage(fmt.Sprintf("%.0f", minTemp), fontType, colornames.Blue, pos, img)

}

func Next90MinTempGraph(img draw.Image, pos image.Point, height, width int) {
	forecast := openweather.NextTwoDaysForecast()
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
	forecasts := openweather.WeekForecast()

	for i, forecast := range forecasts {
		if i >= numDays {
			break
		}
		iconSize := 17
		icon := icons.GetWeatherIcon(forecast.Weather[0].Icon, iconSize)

		draw.Draw(img, image.Rect(pos.X, pos.Y, pos.X+iconSize, pos.Y+iconSize), icon, image.Point{}, draw.Over)
		text.WriteOnImage(fmt.Sprintf("%.0f", forecast.FeelsLike.Day), text.ExtraSmall, color.White, pos.Add(image.Point{iconSize + 3, 4}), img)
		text.WriteOnImage(fmt.Sprintf("%.0f", forecast.FeelsLike.Night), text.ExtraSmall, color.White, pos.Add(image.Point{iconSize + 21, 4}), img)

		pos.Y += iconSize
	}
}
