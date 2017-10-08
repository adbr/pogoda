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
	"strings"
	"text/template"
	"time"

	"github.com/adbr/pogoda/openweather"
)

const (
	serviceURL    = "http://api.openweathermap.org/data/2.5/weather"
	serviceApiKey = "93ca2c840c952abe90064d9e251347f1"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("pogoda: ")

	help := flag.Bool("h", false, "Wyświetla help")
	days := flag.Int("d", 0, "Prognoza pogody na n dni")

	flag.Usage = usage
	flag.Parse()

	if *help {
		fmt.Print(helpStr)
		os.Exit(0)
	}

	if flag.NArg() < 1 {
		usage()
		os.Exit(1)
	}
	city := flag.Arg(0)

	if *days > 0 {
		// Prognoza pogody na *days dni
		forecast, err := openweather.GetForecast(city, *days)
		if err != nil {
			log.Fatal(err)
		}
		err = openweather.PrintForecast(os.Stdout, forecast)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// Pogoda na teraz
		weather, err := getWeather(city)
		if err != nil {
			log.Fatal(err)
		}
		err = printWeather(os.Stdout, weather)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// usage drukuje na stderr sposób użycia programu.
func usage() {
	const s = "Sposób użycia: pogoda [-h] [-d days] miasto"
	fmt.Fprintf(os.Stderr, "%s\n", s)
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

	// Formatowanie danych lokalnych dla ułatwienia wyświetlania
	// niektórych pól.
	l := "15:04:05 MST" // format czasu
	result.Local.SunriseTime = time.Unix(result.Sys.SunriseUnix, 0).Format(l)
	result.Local.SunsetTime = time.Unix(result.Sys.SunsetUnix, 0).Format(l)
	result.Local.Weather = weatherDescription(result.Weather)

	return result, nil
}

// printWeather drukuje do out dane pogodowe sformatowane przy użyciu
// template'u.
func printWeather(out io.Writer, weather *WeatherResult) error {
	err := templ.Execute(out, weather)
	if err != nil {
		return err
	}
	return nil
}

// Typ WeatherResult representuje dane pogodowe.
type WeatherResult struct {
	Coord struct {
		Lat float64 // city geo location, latitude
		Lon float64 // city geo location, longitude
	}
	Weather Weathers
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
		Type        int
		Id          int
		Message     float64
		Country     string
		SunriseUnix int64 `json:"sunrise"`
		SunsetUnix  int64 `json:"sunset"`
	}
	Id   int    // city id
	Name string // city name
	Cod  int

	// Pole Local zawiera dane dodane lokalnie dla ułatwienia
	// wyświetlania informacji.
	Local struct {
		SunriseTime string // sformatowany czas z SunriseUnix
		SunsetTime  string // sformatowany czas z SunsetUnix
		Weather     string // sformatowany opis pogody z pola Weather
	}
}

// Weathers jest typem pola WeatherResult.Weather. Zawiera słowne
// opisy pogody.
type Weathers []struct {
	Id          int
	Main        string
	Description string
	Icon        string
}

// weatherDescription zwraca sformatowany opis pogody.
func weatherDescription(w Weathers) string {
	var a []string
	for _, d := range w {
		s := fmt.Sprintf("%s (%s)", d.Main, d.Description)
		a = append(a, s)
	}
	return strings.Join(a, ", ")
}

// templStr jest templatem dla wyświetlania danych pogodowych typu
// WeatherResult.
const templStr = `Miasto:	       {{.Name}}, {{.Sys.Country}} [{{.Coord.Lat}}, {{.Coord.Lon}}]
Temperatura:   {{.Main.Temp}} °C
{{- if ne .Main.TempMin .Main.TempMax}} (min: {{.Main.TempMin}}, max: {{.Main.TempMax}})
{{- end}}
Pogoda:        {{.Local.Weather}}
Zachmurzenie:  {{.Clouds.All}} %
Wiatr:         {{.Wind.Speed}} m/s ({{.Wind.Deg}}°)
Ciśnienie:     {{.Main.Pressure}} hPa
Wilgotność:    {{.Main.Humidity}} %
Wschód słońca: {{.Local.SunriseTime}}
Zachód słońca: {{.Local.SunsetTime}}
`

var templ = template.Must(template.New("weather").Parse(templStr))

const helpStr = `Program pogoda wyświetla dane pogodowe dla podanego miasta.

Sposób użycia:
	pogoda [-h] [-d days] miasto

	-h	wyświetla help
	-d days
		prognoza pogody na days dni (days: [1..17])

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
