package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

// Chapter type
type Chapter struct {
	Title   string   `json:"title"`
	Content []string `json:"story"`
	Options []Option `json:"options"`
}

// Option type
type Option struct {
	Description string `json:"text"`
	Arc         string `json:"arc"`
}

// ServerHandleFunc type
type ServerHandleFunc struct {
	route   string
	content Chapter
}

func (srv *ServerHandleFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl := template.Must(template.ParseFiles("templates/chapter.html"))
	tmpl.Execute(w, srv.content)
}

func main() {
	story := loadStory()
	for k, chapter := range story {
		http.Handle("/"+k, &ServerHandleFunc{k, chapter})
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/intro", 301)
	})
	fmt.Println("Listen on port:8080 ....")
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func loadStory() map[string]Chapter {
	content, err := ioutil.ReadFile("gopher.json")
	check(err)
	var story map[string]Chapter
	err = json.Unmarshal(content, &story)
	check(err)
	return story
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
