package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Story struct {
	Title   string   `json:"title"`
	Text    []string `json:"story"`
	Options []Option `json:"options"`
}

type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

var arcs = map[string]Story{}

func init() {
	f, _ := os.Open("gophers.json")
	defer f.Close()
	bytes, _ := ioutil.ReadAll(f)
	err := json.Unmarshal(bytes, &arcs)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Gopher Adventure Json parsed")
}

func main() {
	log.Println("Waiting")
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path[1:]
	arc, err := getStory(urlPath)
	if err != nil {
		t, innerErr := template.ParseFiles("error.html")
		if innerErr != nil {
			log.Fatal(innerErr.Error())
		}
		innerErr = t.Execute(w, nil)
		if innerErr != nil {
			log.Fatal(err.Error())
		}
		log.Printf("Rendered error page, story arc '%s' does not exist", urlPath)

	} else {

		t, err := template.ParseFiles("index.html")
		if err != nil {
			log.Fatal(err.Error())
		} else {
			data := make(map[string]interface{})
			data["arc"] = arc
			err := t.Execute(w, data)
			if err != nil {
				log.Fatal(err.Error())
			} else {
				log.Printf("Rendered story arc '%s'", arc.Title)
			}
		}
	}
}

func getStory(title string) (Story, error) {
	for k, v := range arcs {
		if k == title {
			return v, nil
		}
	}
	return Story{}, fmt.Errorf("Title %s does not exist", title)
}
