// 2017-02-04 adbr

// TODO: pakiet dla openweathermap?

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"text/template"
	"time"
)

const (
	serviceURL    = "http://api.openweathermap.org/data/2.5/weather"
	serviceApiKey = "93ca2c840c952abe90064d9e251347f1"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("pogoda: ")

	flag.Usage = usage
	flag.Parse()

	if flag.NArg() < 1 {
		usage()
		os.Exit(1)
	}
	city := flag.Arg(0)

	weather, err := getWeather(city)
	if err != nil {
		log.Fatal(err)
	}

	err = printWeather(os.Stdout, weather)
	if err != nil {
		log.Fatal(err)
	}
}

const usageStr = `Sposób użycia: pogoda miasto
`

// usage drukuje na stderr sposób użycia programu.
func usage() {
	fmt.Fprint(os.Stderr, usageStr)
}

// Typ WeatherResult representuje dane pogodowe zdekodowane z JSONa.
type WeatherResult struct {
	Coord struct {
		Lat float64 // city geo location, latitude
		Lon float64 // city geo location, longitude
	}
	Weather []struct {
		Id          int
		Main        string
		Description string
		Icon        string
	}
	Base string
	Main struct {
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
		Type        int
		Id          int
		Message     float64
		Country     string
		SunriseUnix int64  `json:"sunrise"`
		SunsetUnix  int64  `json:"sunset"`
		SunriseTime string // sformatowany czas z SunriseUnix
		SunsetTime  string // sformatowany czas z SunsetUnix
	}
	Id   int    // city id
	Name string // city name
	Cod  int
}

// getWeather zwraca pogodę dla miasta city.
func getWeather(city string) (*WeatherResult, error) {
	query := url.Values{
		"appid": {serviceApiKey},
		"units": {"metric"},
		"q":     {city},
	}
	resp, err := http.Get(serviceURL + "?" + query.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Wczytanie zwróconych danych w formacie JSON.
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Dekodowanie JSONa.
	result := new(WeatherResult)
	err = json.Unmarshal(data, result)
	if err != nil {
		return nil, err
	}
	l := "15:04:05 MST"
	result.Sys.SunriseTime = time.Unix(result.Sys.SunriseUnix, 0).Format(l)
	result.Sys.SunsetTime = time.Unix(result.Sys.SunsetUnix, 0).Format(l)
	return result, nil
}

const templStr = `Miasto:	       {{.Name}}, {{.Sys.Country}} [{{.Coord.Lat}}, {{.Coord.Lon}}]
Temperatura:   {{.Main.Temp}} °C (min: {{.Main.TempMin}}, max: {{.Main.TempMax}})
{{range .Weather -}}
Pogoda:        {{.Main}} ({{.Description}})
{{- end}}
Ciśnienie:     {{.Main.Pressure}} hpa
Wilgotność:    {{.Main.Humidity}} %
Wiatr:         {{.Wind.Speed}} m/s ({{.Wind.Deg}}°)
Zachmurzenie:  {{.Clouds.All}} %
Wschód słońca: {{.Sys.SunriseTime}}
Zachód słońca: {{.Sys.SunsetTime}}
(Dane pochodzą z serwisu OpenWeatherMap.com)
`

var templ = template.Must(template.New("weather").Parse(templStr))

func printWeather(out io.Writer, weather *WeatherResult) error {
	err := templ.Execute(out, weather)
	if err != nil {
		return err
	}
	return nil
}
