// 2017-02-04 adbr

// TODO: pakiet dla openweathermap?

package main

import (
	"encoding/json"
	"flag"
	"fmt"
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
		log.Print("brakuje nazwy miasta")
		usage()
		os.Exit(1)
	}
	city := flag.Arg(0)

	weather, err := getWeather(city)
	if err != nil {
		log.Fatal(err)
	}

	err = printWeather(weather)
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
	Name string
	Main struct {
		Temp     float64
		Pressure float64
		Humidity float64
		TempMin  float64 `json:"temp_min"`
		TempMax  float64 `json:"temp_max"`
	}
	Wind struct {
		Speed float64
	}
	Clouds struct {
		All float64
	}
	Sys struct {
		Country     string
		SunriseUnix int64 `json:"sunrise"`
		SunsetUnix  int64 `json:"sunset"`
		Sunrise     time.Time
		Sunset      time.Time
	}
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
	result.Sys.Sunrise = time.Unix(result.Sys.SunriseUnix, 0)
	result.Sys.Sunset = time.Unix(result.Sys.SunsetUnix, 0)
	return result, nil
}

const templStr = `Miasto:	       {{.Name}}, {{.Sys.Country}}
Temperatura:   {{.Main.Temp}} °C (min: {{.Main.TempMin}}, max: {{.Main.TempMax}})
Ciśnienie:     {{.Main.Pressure}} hpa
Wilgotność:    {{.Main.Humidity}} %
Wiatr:         {{.Wind.Speed}} m/s
Zachmurzenie:  {{.Clouds.All}} %
Wschód słońca: {{.Sys.Sunrise}}
Zachód słońca: {{.Sys.Sunset}}
[Dane pochodzą z serwisu OpenWeatherMap]
`

var templ = template.Must(template.New("weather").Parse(templStr))

func printWeather(weather *WeatherResult) error {
	err := templ.Execute(os.Stdout, weather)
	if err != nil {
		return err
	}
	return nil
}
