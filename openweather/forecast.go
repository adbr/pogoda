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
	"time"
	"strings"
)

const (
	forecastURL   = "http://api.openweathermap.org/data/2.5/forecast/daily"
	serviceApiKey = "93ca2c840c952abe90064d9e251347f1"
)

// GetForecast pobiera z serwisu i zwraca prognozę pogody dla miasta
// city na days dni naprzód. Liczba dni days musi być z przedziału
// [1..16]. Zwraca ForecastResult i błąd jeśli wystąpił.
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
		return nil, fmt.Errorf("zapytanie %s: %s", query, resp.Status)
	}

	// Wczytanie zwróconych danych.
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("czytanie danych pogodowych: %s", err)
	}

	// Dekodowanie JSONa.
	result := &ForecastResult{}
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

type WeatherDescription struct {
	Id          int
	Main        string
	Description string
	Icon        string
}

// FormatUnixTime zwraca sformatowany czas odpowiadający czasowi w
// formacie Unixa u.
func FormatUnixTime(u int64) string {
	t := time.Unix(u, 0)
	s := fmt.Sprintf("%02d.%02d.%d %s",
		t.Day(), t.Month(), t.Year(), weekdays[t.Weekday()])
	return s
}

// weekdays zawiera polskie nazwy dni tygodnia.
var weekdays = [...]string{
	"Niedziela",
	"Poniedziałek",
	"Wtorek",
	"Środa",
	"Czwartek",
	"Piątek",
	"Sobota",
}

// FormatWeather zwraca sformatowany opis pogody.
func FormatWeather(w []WeatherDescription) string {
	var a []string
	for _, wd := range w {
		s := fmt.Sprintf("%s (%s)", wd.Main, wd.Description)
		a = append(a, s)
	}
	return strings.Join(a, ", ")
}

// funcMap jest mapą funkcji dla template.
var funcMap = template.FuncMap{
	"FormatUnixTime": FormatUnixTime,
	"FormatWeather": FormatWeather,
}

// forecastTmplText jest templatem dla wyświetlania danych pogodowych
// typu ForecastResult.
const forecastTmplText = `Miasto: {{.City.Name}}, {{.City.Country}} [{{.City.Coord.Lat}}, {{.City.Coord.Lon}}]
Prognoza na {{.Cnt}} dni
{{range .List -}}
========================================
Dzień: {{FormatUnixTime .Dt}}
Temperatura:
    Dzień:     {{.Temp.Day}} °C (Min: {{.Temp.Min}}, Max: {{.Temp.Max}})
    Noc:       {{.Temp.Night}} °C
    Wieczór:   {{.Temp.Eve}} °C
    Rano:      {{.Temp.Morn}} °C
Pogoda:        {{FormatWeather .Weather}}
Zachmurzenie:  {{.Clouds}} %
Wiatr:         {{.Speed}} m/s ({{.Deg}}°)
Ciśnienie:     {{.Pressure}} hPa
Wilgotność:    {{.Humidity}} %
{{end}}`

var forecastTmpl = template.Must(
	template.New("forecast").Funcs(funcMap).Parse(forecastTmplText))
