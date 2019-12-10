package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type apiData struct {
	ParamName string
	Key       string
}

func (api apiData) getAPIParam() string {
	return api.ParamName + "=" + api.Key
}

func serveReverseProxy(target string, apiKey apiData, res http.ResponseWriter, req *http.Request) {
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

	fmt.Println(req.URL.String())
	fmt.Println(req.URL.Path)
	fmt.Println(req.URL.RawQuery)

	// Note that ServeHttp is non blocking and uses a go routine under the hood
	proxy.ServeHTTP(res, req)
}

func handleAutocomplete(res http.ResponseWriter, req *http.Request) {
	targetURL := "https://api.weatherstack.com"

	keyData := apiData{ParamName: "access_key", Key: "00a232561814c1aa1784c70836d537d9"}

	serveReverseProxy(targetURL, keyData, res, req)

}

func main() {

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
	http.HandleFunc("/api/autocomplete", handleAutocomplete)

	http.ListenAndServe(addr, nil)

}
