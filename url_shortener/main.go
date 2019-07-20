package main

import (
	"flag"
	"fmt"
	"gophercises-training/url_shortener/urlshort"
	"io/ioutil"
	"net/http"
)

func main() {
	mux := defaultMux()

	fileYAML := flag.String("yaml", "redirects.yaml", "The YAML file with the redirects config")
	fileJSON := flag.String("json", "redirects.json", "The JSON file with the redirects config")
	flag.Parse()

	redirectsYAML, err := ioutil.ReadFile(*fileYAML)
	check(err)

	redirectsJSON, err := ioutil.ReadFile(*fileJSON)
	check(err)

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yaml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}

	yamlFileHandler, err := urlshort.YAMLHandler([]byte(string(redirectsYAML)), yamlHandler)
	if err != nil {
		panic(err)
	}

	jsonFileHandler, err := urlshort.JSONHandler([]byte(string(redirectsJSON)), yamlFileHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonFileHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
