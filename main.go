// 2017-02-04 adbr

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const (
	serviceURL    = "http://api.openweathermap.org/data/2.5/weather"
	serviceApiKey = "93ca2c840c952abe90064d9e251347f1"
)

func main() {
	city := "Krakow,pl"

	// Utworzenie URLa z zapytaniem
	query := url.Values{
		"q":     {city},
		"appid": {serviceApiKey},
	}
	urlStr := serviceURL + "?" + query.Encode()

	log.Printf("url: %s\n", urlStr) // debug

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
