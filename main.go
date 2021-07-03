package main

import (
	"encoding/json"
	"fmt"
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
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	arc, err := getStory(r.URL.Path[1:])
	if err != nil {
		fmt.Fprintf(w, fmt.Sprintf("<html><body>%v</body></html>", err.Error()))
	} else {
		fmt.Fprintf(w, fmt.Sprintf("<html><body>title:%s<br>story:%v<br>options:%v</body></html>", arc.Title, arc.Text, arc.Options))
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
