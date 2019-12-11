package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

var yaAPIKEey string
var apixuAPIKEey string 

type apiData struct {
	ParamName string
	Key       string
}

func (api apiData) getAPIParam() string {
	return api.ParamName + "=" + api.Key
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
}

func serveWeatherstackSP(target string, apiKey apiData, res http.ResponseWriter, req *http.Request) {
	// parse the url
	url, _ := url.Parse(target)

	targetQuery := req.URL.RawQuery + "&" + apiKey.getAPIParam()

	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)

	// Update the headers to allow for SSL redirection
	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme

	req.Host = url.Host
	req.URL.Path = "autocomplete"
	req.URL.RawQuery = targetQuery

	proxy.ServeHTTP(res, req)
}

func serveYandexWeatherSP(target string, apiKey string, res http.ResponseWriter, req *http.Request) {
	// parse the url
	url, _ := url.Parse(target)

	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)

	// Update the headers to allow for SSL redirection
	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme

	req.Host = url.Host
	req.URL.Path = "v1/forecast"

	req.Header.Set("X-Yandex-API-Key", apiKey)

	enableCors(&res)

	proxy.ServeHTTP(res, req)
}

func handleAutocomplete(res http.ResponseWriter, req *http.Request) {

	targetURL := "https://api.weatherstack.com"

	

	keyData := apiData{ParamName: "access_key", Key: apixuAPIKEey}

	serveWeatherstackSP(targetURL, keyData, res, req)

}

func handleForecast(res http.ResponseWriter, req *http.Request) {
	
	targetURL := "https://api.weather.yandex.ru"

	serveYandexWeatherSP(targetURL, yaAPIKEey, res, req)

}

func main() {

	_ = godotenv.Load()

	apixuAPIKEey = os.Getenv("WEATHERSTACK_API_KEY")
	yaAPIKEey = os.Getenv("YA_WEATHER_API_KEY")
	
	if apixuAPIKEey == "" || yaAPIKEey == "" {
		log.Fatal("Failed to initialize API variables")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		query, ok := req.URL.Query()["query"]
		if ok {
			fmt.Fprintf(w, "Hello, %q", query[0])
		}

	})

	port := os.Getenv("PORT")

	if port == "" {
		port = "9000"
	}
	addr := ":" + port
	log.Printf("Listening on %s...\n", addr)
	http.HandleFunc("/api/autocomplete/", handleAutocomplete)
	http.HandleFunc("/api/forecast/", handleForecast)
	http.ListenAndServe(addr, nil)

}
