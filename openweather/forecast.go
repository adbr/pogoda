// 2017-09-28 adbr

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
	forecastURL = "http://api.openweathermap.org/data/2.5/forecast/daily"
)

// GetForecast pobiera z serwisu i zwraca prognozę pogody dla miasta
// city na days dni naprzód. Liczba dni days musi być z przedziału
// [1..17]. Zwraca ForecastResult i błąd jeśli wystąpił.
func GetForecast(city string, days int) (*ForecastResult, error) {
	queryValues := url.Values{
		"appid": {serviceApiKey},
		"units": {"metric"},
		"q":     {city},
		"cnt":   {fmt.Sprintf("%d", days)},
	}

	// Zapytanie do serwisu.
	query := forecastURL + "?" + queryValues.Encode()
	resp, err := http.Get(query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("zapytanie %s: %s: %s",
			forecastURL, city, resp.Status)
	}

	// Wczytanie zwróconych danych.
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("czytanie danych pogodowych: %s", err)
	}

	// Dekodowanie JSONa.
	result := new(ForecastResult)
	err = json.Unmarshal(data, result)
	if err != nil {
		return nil, fmt.Errorf("dekodowanie danych pogodowych: %s", err)
	}

	return result, nil
}

// PrintForecast drukuje do out sformatowane dane pogodowe prognozy.
func PrintForecast(out io.Writer, forecast *ForecastResult) error {
	err := forecastTmpl.Execute(out, forecast)
	if err != nil {
		return fmt.Errorf("formatowanie danych pogodowych: %s", err)
	}
	return nil
}

// ForecastResult representuje dane pogodowe prognozy.
type ForecastResult struct {
	City struct {
		Id    int
		Name  string
		Coord struct {
			Lon float64
			Lat float64
		}
		Country    string
		Population int
	}
	Cod     string
	Message float64
	Cnt     int
	List    []DayWeather
}

// DayWeather jest elementem typu ForecastResult i zawiera dane
// prognozy pogody dla pojedynczego dnia.
type DayWeather struct {
	Dt   int64
	Temp struct {
		Day   float64
		Min   float64
		Max   float64
		Night float64
		Eve   float64
		Morn  float64
	}
	Pressure float64
	Humidity float64
	Weather  []WeatherDescription
	Speed    float64
	Deg      float64
	Clouds   float64
	Rain     float64
}

// WeatherDescription jest elementem typów DayWeather i WeatherResult
// i zawiera słowny opis pogody.
type WeatherDescription struct {
	Id          int
	Main        string
	Description string
	Icon        string
}

// forecastTmplText jest templatem dla wyświetlania danych pogodowych
// typu ForecastResult.
const forecastTmplText = `Miasto: {{.City.Name}}, {{.City.Country}} [{{.City.Coord.Lat}}, {{.City.Coord.Lon}}]
Prognoza na {{.Cnt}} dni:
{{range .List -}}
========================================
Dzień:         {{formatUnixDate .Dt}}
Temperatura:
    Dzień:     {{.Temp.Day}} °C (Min: {{.Temp.Min}}, Max: {{.Temp.Max}})
    Noc:       {{.Temp.Night}} °C
    Wieczór:   {{.Temp.Eve}} °C
    Rano:      {{.Temp.Morn}} °C
Pogoda:        {{formatWeather .Weather}}
Zachmurzenie:  {{.Clouds}} %
Wiatr:         {{.Speed}} m/s ({{.Deg}}°)
Ciśnienie:     {{.Pressure}} hPa
Wilgotność:    {{.Humidity}} %
{{end}}`

var forecastTmpl = template.Must(
	template.New("forecast").Funcs(funcMap).Parse(forecastTmplText))
