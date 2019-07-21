package main

import (
	"flag"
	"fmt"
	"gophercises-training/url_shortener/urlshort"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/boltdb/bolt"
)

func main() {
	mux := defaultMux()

	populateDB()

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

	dbHandler, err := urlshort.BoltDBHandler(getRedirectsDB(), jsonFileHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", dbHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func populateDB() {
	db, err := bolt.Open("redirects.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Redirects"))
		if err != nil {
			return fmt.Errorf("Create bucket: %s", err)
		}
		return nil
	})

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Redirects"))
		err := b.Put([]byte("/urlshort-4"), []byte("https://github.com/gophercises/urlshort"))
		return err
	})

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Redirects"))
		err := b.Put([]byte("/urlshort-final-4"), []byte("https://github.com/gophercises/urlshort/tree/solution"))
		return err
	})
}

func getRedirectsDB() map[string]string {
	db, err := bolt.Open("redirects.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	pathsToUrls := map[string]string{}

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Redirects"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("key=%s, value=%s\n", k, v)
			pathsToUrls[string(k)] = string(v)
		}
		return nil
	})
	return pathsToUrls
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
