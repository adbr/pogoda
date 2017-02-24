// 2017-02-04 adbr

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
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

	// Utworzenie URLa z zapytaniem
	query := url.Values{
		"appid": {serviceApiKey},
		"q":     {city},
	}
	urlStr := serviceURL + "?" + query.Encode()

	fmt.Printf("[debug] url: %s\n", urlStr) // debug

	// Request HTTP
	resp, err := http.Get(urlStr)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Wczytanie zwróconych danych
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Wydruk zwróconych danych
	fmt.Printf("%s\n", data)
}

// usage drukuje na stderr sposób użycia programu.
func usage() {
	fmt.Fprint(os.Stderr, usageStr)
}

const usageStr = `Sposób użycia: pogoda miasto
`
