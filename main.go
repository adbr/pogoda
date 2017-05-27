// 2017-02-04 adbr

// Program pogoda wyświetla dane pogodowe dla podanego miasta.
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

	h := flag.Bool("h", false, "Wyświetla help")

	flag.Usage = usage
	flag.Parse()

	if *h {
		fmt.Print(helpStr)
		os.Exit(0)
	}

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

// usage drukuje na stderr sposób użycia programu.
func usage() {
	const s = "Sposób użycia: pogoda [-h] miasto"
	fmt.Fprintf(os.Stderr, "%s\n", s)
}

const helpStr = `Program pogoda wyświetla dane pogodowe dla podanego miasta.

Sposób użycia:
	pogoda [-h] miasto

	-h Wyświetla help.

Dla podanego miasta program pobiera aktualne dane pogodowe z serwisu
http://api.openweathermap.org i wyświetla je na standardowe wyjście.

Przykład: pogoda dla Warszawy:

	$ pogoda Warszawa
	Miasto:        Warszawa, PL [52.24, 21.04]
	Temperatura:   21 °C (min: 21, max: 21)
	Pogoda:        Clear (clear sky)
	Ciśnienie:     1023 hpa
	Wilgotność:    40 %
	Wiatr:         2.6 m/s (0°)
	Zachmurzenie:  0 %
	Wschód słońca: 04:24:40 CEST
	Zachód słońca: 20:42:16 CEST
	(Dane pochodzą z serwisu OpenWeatherMap.com)
`

// Typ WeatherResult representuje dane pogodowe.
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

// getWeather zwraca dane pogodowe dla miasta city. Dane są pobierane
// z serwisu openweathermap.org.
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
	l := "15:04:05 MST" // format czasu
	result.Sys.SunriseTime = time.Unix(result.Sys.SunriseUnix, 0).Format(l)
	result.Sys.SunsetTime = time.Unix(result.Sys.SunsetUnix, 0).Format(l)

	return result, nil
}

// templStr jest templatem dla wyświetlania danych pogodowych typu
// WeatherResult.
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

// printWeather drukuje do out dane pogodowe sformatowane przy użyciu
// template'u.
func printWeather(out io.Writer, weather *WeatherResult) error {
	err := templ.Execute(out, weather)
	if err != nil {
		return err
	}
	return nil
}
