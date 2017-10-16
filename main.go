// 2017-02-04 adbr

/*
Program pogoda wyświetla dane pogodowe dla podanego miasta.

Sposób użycia:
	pogoda [opcje] miasto
	-h      sposób użycia
	-help   dokumentacja
	-d days
	        prognoza pogody na days dni (days: [1..17])

Dla podanego miasta program pobiera aktualne dane pogodowe z serwisu
http://api.openweathermap.org i wyświetla je na standardowe wyjście. Z
opcją -d wyświetla prognozę pogody na maksymalnie 16 dni.

Przykład: pogoda dla Warszawy:

	$ pogoda Warszawa
	Miasto:        Warszawa, PL [52.24, 21.04]
	Data:          13.10.2017 13:00:00 CEST, Piątek
	Temperatura:   13 °C
	Pogoda:        Clouds (broken clouds)
	Zachmurzenie:  75 %
	Wiatr:         7.7 m/s (290°)
	Ciśnienie:     1017 hPa
	Wilgotność:    66 %
	Wschód słońca: 06:58:17 CEST
	Zachód słońca: 17:44:40 CEST
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/adbr/pogoda/openweather"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("pogoda: ")

	help := flag.Bool("h", false, "wyświetla help")
	helpLong := flag.Bool("help", false, "wyświetla dokumentację")
	days := flag.Int("d", 0, "prognoza pogody na n dni")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usageText)
	}
	flag.Parse()

	if *help {
		fmt.Print(usageText)
		os.Exit(0)
	}
	if *helpLong {
		fmt.Print(helpText)
		os.Exit(0)
	}

	if flag.NArg() < 1 {
		fmt.Fprint(os.Stderr, usageText)
		os.Exit(1)
	}
	city := flag.Arg(0)

	if *days != 0 {
		// Prognoza pogody na *days dni
		if *days < 1 || *days > 17 {
			log.Printf("zła wartość opcji -d: %d\n", *days)
			fmt.Fprint(os.Stderr, usageText)
			os.Exit(1)
		}
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
		weather, err := openweather.GetWeather(city)
		if err != nil {
			log.Fatal(err)
		}
		err = openweather.PrintWeather(os.Stdout, weather)
		if err != nil {
			log.Fatal(err)
		}
	}
}

const usageText = `Sposób użycia:
	pogoda [opcje] miasto
	-h	sposób użycia
	-help	dokumentacja
	-d days
		prognoza pogody na days dni (days: [1..17])
`

const helpText = `Program pogoda wyświetla dane pogodowe dla podanego miasta.

Sposób użycia:
	pogoda [opcje] miasto
	-h      sposób użycia
	-help   dokumentacja
	-d days
	        prognoza pogody na days dni (days: [1..17])

Dla podanego miasta program pobiera aktualne dane pogodowe z serwisu
http://api.openweathermap.org i wyświetla je na standardowe wyjście. Z
opcją -d wyświetla prognozę pogody na maksymalnie 16 dni.

Przykład: pogoda dla Warszawy:

	$ pogoda Warszawa
	Miasto:        Warszawa, PL [52.24, 21.04]
	Data:          13.10.2017 13:00:00 CEST, Piątek
	Temperatura:   13 °C
	Pogoda:        Clouds (broken clouds)
	Zachmurzenie:  75 %
	Wiatr:         7.7 m/s (290°)
	Ciśnienie:     1017 hPa
	Wilgotność:    66 %
	Wschód słońca: 06:58:17 CEST
	Zachód słońca: 17:44:40 CEST
`
