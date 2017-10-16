// 2017-10-16 adbr

// Pakiet openweather udostępnia funkcjonalność związaną z pobieraniem
// i drukowaniem danych pogodowych z serwisu
// http://openweathermap.org/.  Umożliwia pobieranie aktualnej pogody
// i prognozy pogody na kilka dni (maksymalnie 16).
package openweather

import (
	"fmt"
	"strings"
	"text/template"
	"time"
)

const (
	serviceApiKey = "93ca2c840c952abe90064d9e251347f1"
)

// funcMap jest mapą funkcji dla template'ów do formatowania niektórych
// danych.
var funcMap = template.FuncMap{
	"formatUnixDate":     formatUnixDate,
	"formatUnixTime":     formatUnixTime,
	"formatUnixDateTime": formatUnixDateTime,
	"formatWeather":      formatWeather,
}

// formatUnixDate zwraca sformatowaną datę odpowiadającą czasowi w
// formacie Unixa.
func formatUnixDate(u int64) string {
	t := time.Unix(u, 0)
	s := fmt.Sprintf("%02d.%02d.%d, %s",
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

// formatUnixTime zwraca sformatowany czas odpowiadający czasowi w
// formacie Unixa.
func formatUnixTime(u int64) string {
	t := time.Unix(u, 0)
	zone, _ := t.Zone()
	s := fmt.Sprintf("%02d:%02d:%02d %s",
		t.Hour(), t.Minute(), t.Second(), zone)
	return s
}

// formatUnixDateTime zwraca sformatowaną datę i czas odpowiadający
// czasowi w formacie Unixa.
func formatUnixDateTime(u int64) string {
	t := time.Unix(u, 0)
	zone, _ := t.Zone()
	s := fmt.Sprintf("%02d.%02d.%d %02d:%02d:%02d %s, %s",
		t.Day(), t.Month(), t.Year(),
		t.Hour(), t.Minute(), t.Second(),
		zone, weekdays[t.Weekday()])
	return s
}

// formatWeather zwraca sformatowany opis pogody.
func formatWeather(w []WeatherDescription) string {
	var a []string
	for _, wd := range w {
		s := fmt.Sprintf("%s (%s)", wd.Main, wd.Description)
		a = append(a, s)
	}
	return strings.Join(a, ", ")
}
