package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	app := app.New().NewWindow("Weather")
	cityList := [][2]string{
		{"Delhi", "delhi"},
		{"Noida", "noida"},
		{"Mumbai", "mumbai"}}
	app.SetContent(makeListTab(cityList, app))
	app.Resize(fyne.Size{Height: 400, Width: 580})
	app.ShowAndRun()
}

type WeatherInfo struct {
	cordinates     [2]float64
	temperature    float64
	minTemperature float64
	maxTemperature float64
	humidity       float64
	pressure       int64
	visibility     float64
	windspeed      float64
	winddir        string
}

func DegToCard(d float64) string {
	dirs := []string{"N", "NNE", "NE", "ENE", "E", "ESE", "SE", "SSE", "S", "SSW", "SW", "WSW", "W", "WNW", "NW", "NNW"}
	return dirs[int64((d+11.25)/22.5)%16]
}

func WeatherInfoFilter(city string) WeatherInfo {
	response, err := http.Get("https://api.openweathermap.org/data/2.5/weather?q=" + city + "&APPID=8bd4d6a9a95aadd819800ea99b951a8e")
	if err != nil {
		log.Panicln(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Panicln(err)
	}
	weather, err := UnmarshalWelcome(body)
	if err != nil {
		log.Panicln(err)
	}
	var info WeatherInfo
	info.temperature = weather.Main.Temp - 273.15
	info.minTemperature = weather.Main.TempMin - 273.15
	info.maxTemperature = weather.Main.TempMax - 273.15
	info.humidity = float64(weather.Main.Humidity) / 100
	info.pressure = weather.Main.Pressure
	info.cordinates = [2]float64{weather.Coord.Lat, weather.Coord.Lon}
	info.visibility = float64(weather.Visibility / 1000)
	info.windspeed = weather.Wind.Speed
	info.winddir = DegToCard(float64(weather.Wind.Deg))
	return info
}

func SelectCity(callback func(string)) *widget.Select {
	cityList := []string{"Delhi", "Noida", "Mumbai"}
	options := widget.NewSelect(cityList, func(s string) {
		callback(s)
	})
	return options
}

func DisplayLayout(city string) *fyne.Container {
	weather := WeatherInfoFilter(city)
	cordinates := [2]string{fmt.Sprintf("%.2f° N", weather.cordinates[0]), fmt.Sprintf("%.2f° E", weather.cordinates[1])}
	liveTemp := widget.NewLabel(fmt.Sprintf("%.2f°C", weather.temperature))
	maxTemp := widget.NewLabel(fmt.Sprintf("%.2f°C", weather.maxTemperature))
	minTemp := widget.NewLabel(fmt.Sprintf("%.2f°C", weather.minTemperature))
	humid := widget.NewProgressBarWithData(binding.BindFloat(&weather.humidity))
	pressure := widget.NewLabel(fmt.Sprintf("%dhPa", weather.pressure))
	visibility := widget.NewLabel(fmt.Sprintf("%.1fKm", weather.visibility))
	wind := widget.NewLabel(fmt.Sprintf("%0.2fm/s %s", weather.windspeed, weather.winddir))
	return container.NewVBox(
		container.New(
			layout.NewGridLayoutWithColumns(2),
			container.New(
				layout.NewGridLayoutWithRows(9),
				container.NewHBox(widget.NewLabel("Lattitude"), layout.NewSpacer(), widget.NewSeparator()),
				container.NewHBox(widget.NewLabel("Longitude"), layout.NewSpacer(), widget.NewSeparator()),
				container.NewHBox(widget.NewLabel("Temperature"), layout.NewSpacer(), widget.NewSeparator()),
				container.NewHBox(layout.NewSpacer(), container.NewCenter(widget.NewLabel("Min. Temp.")), layout.NewSpacer(), widget.NewSeparator()),
				container.NewHBox(layout.NewSpacer(), container.NewCenter(minTemp), layout.NewSpacer(), widget.NewSeparator()),
				container.NewHBox(widget.NewLabel("Humidity"), layout.NewSpacer(), widget.NewSeparator()),
				container.NewHBox(widget.NewLabel("Pressure"), layout.NewSpacer(), widget.NewSeparator()),
				container.NewHBox(widget.NewLabel("Visibility"), layout.NewSpacer(), widget.NewSeparator()),
				container.NewHBox(widget.NewLabel("Wind"), layout.NewSpacer(), widget.NewSeparator()),
			),
			container.New(
				layout.NewGridLayoutWithRows(9),
				container.NewCenter(widget.NewLabel(cordinates[0])),
				container.NewCenter(widget.NewLabel(cordinates[1])),
				container.NewCenter(liveTemp),
				container.NewCenter(widget.NewLabel("Max. Temp.")),
				container.NewCenter(container.NewCenter(maxTemp)),
				container.New(layout.NewGridLayoutWithColumns(1), humid),
				container.NewCenter(pressure),
				container.NewCenter(visibility),
				container.NewCenter(wind),
			),
		),
	)
}

func UnmarshalWelcome(data []byte) (Welcome, error) {
	var r Welcome
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Welcome) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Welcome struct {
	Coord      Coord     `json:"coord"`
	Weather    []Weather `json:"weather"`
	Base       string    `json:"base"`
	Main       Main      `json:"main"`
	Visibility int64     `json:"visibility"`
	Wind       Wind      `json:"wind"`
	Clouds     Clouds    `json:"clouds"`
	Dt         int64     `json:"dt"`
	Sys        Sys       `json:"sys"`
	Timezone   int64     `json:"timezone"`
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Cod        int64     `json:"cod"`
}

type Clouds struct {
	All int64 `json:"all"`
}

type Coord struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

type Main struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	Pressure  int64   `json:"pressure"`
	Humidity  int64   `json:"humidity"`
}

type Sys struct {
	Type    int64  `json:"type"`
	ID      int64  `json:"id"`
	Country string `json:"country"`
	Sunrise int64  `json:"sunrise"`
	Sunset  int64  `json:"sunset"`
}

type Weather struct {
	ID          int64  `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type Wind struct {
	Speed float64 `json:"speed"`
	Deg   int64   `json:"deg"`
}

func makeListTab(data [][2]string, app fyne.Window) fyne.CanvasObject {
	icon := widget.NewIcon(theme.MenuIcon())
	label := widget.NewLabel("Select a city from the menu")
	hbox := container.NewVBox(container.NewHBox(icon, label), layout.NewSpacer())

	list := widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewIcon(theme.MoreVerticalIcon()), widget.NewLabel("Template Object"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(data[id][0])
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		label.SetText(data[id][0])
		icon.SetResource(theme.NavigateNextIcon())
		hbox.Objects = hbox.Objects[0:1]
		hbox.AddObject(DisplayLayout(data[id][1]))
		app.Content().Refresh()
	}
	list.OnUnselected = func(id widget.ListItemID) {
		label.SetText("Select An Item From The List")
		icon.SetResource(nil)
	}

	r := container.NewHSplit(list, hbox)
	r.Offset = 0
	return r
}
