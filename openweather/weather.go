// 2017-10-11 adbr

package openweather

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"text/template"
)

const (
	weatherURL = "http://api.openweathermap.org/data/2.5/weather"
)

// GetWeather pobiera z serwisu i zwraca aktualne dane pogodowe dla
// miasta city. Zwraca WeatherResult i błąd jeśli wystąpił.
func GetWeather(city string) (*WeatherResult, error) {
	queryValues := url.Values{
		"appid": {serviceApiKey},
		"units": {"metric"},
		"q":     {city},
	}

	// Zapytanie do serwisu.
	query := weatherURL + "?" + queryValues.Encode()
	resp, err := http.Get(query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("zapytanie %s: %s: %s",
			weatherURL, city, resp.Status)
	}

	// Wczytanie zwróconych danych.
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("czytanie danych pogodowych: %s", err)
	}

	// Dekodowanie JSONa.
	result := new(WeatherResult)
	err = json.Unmarshal(data, result)
	if err != nil {
		return nil, fmt.Errorf("dekodowanie danych pogodowych: %s", err)
	}

	return result, nil
}

// PrintWeather drukuje do out sformatowane dane pogodowe.
func PrintWeather(out io.Writer, weather *WeatherResult) error {
	err := weatherTmpl.Execute(out, weather)
	if err != nil {
		return fmt.Errorf("formatowanie danych pogodowych: %s", err)
	}
	return nil
}

// WeatherResult representuje aktualne dane pogodowe.
type WeatherResult struct {
	Coord struct {
		Lat float64
		Lon float64
	}
	Weather []WeatherDescription
	Base    string
	Main    struct {
		Temp     float64
		Pressure float64
		Humidity float64
		TempMin  float64 `json:"temp_min"`
		TempMax  float64 `json:"temp_max"`
	}
	Visibility int
	Wind       struct {
		Speed float64
		Deg   float64
	}
	Clouds struct {
		All float64
	}
	Dt  int64
	Sys struct {
		Type    int
		Id      int
		Message float64
		Country string
		Sunrise int64
		Sunset  int64
	}
	Id   int
	Name string
	Cod  int
}

// weatherTmplText jest templatem dla wyświetlania danych pogodowych
// typu WeatherResult.
const weatherTmplText = `Miasto:	       {{.Name}}, {{.Sys.Country}} [{{.Coord.Lat}}, {{.Coord.Lon}}]
Czas:          {{formatUnixDateTime .Dt}}
Temperatura:   {{.Main.Temp}} °C
{{- if ne .Main.TempMin .Main.TempMax}} (min: {{.Main.TempMin}}, max: {{.Main.TempMax}})
{{- end}}
Pogoda:        {{formatWeather .Weather}}
Zachmurzenie:  {{.Clouds.All}} %
Wiatr:         {{.Wind.Speed}} m/s ({{.Wind.Deg}}°)
Ciśnienie:     {{.Main.Pressure}} hPa
Wilgotność:    {{.Main.Humidity}} %
Wschód słońca: {{formatUnixTime .Sys.Sunrise}}
Zachód słońca: {{formatUnixTime .Sys.Sunset}}
`

var weatherTmpl = template.Must(
	template.New("weather").Funcs(funcMap).Parse(weatherTmplText))
