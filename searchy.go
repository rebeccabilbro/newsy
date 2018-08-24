package main

import (
	"fmt"
	"net/http"
	"strings"
)

type Story struct {
	title  string
	url    string
	author string
	source string
}

var stories []Story

func init() {
	stories = append(stories,
		Story{
			"The Actor Model for Actors and Models",
			"https://rebeccabilbro.github.io/actor-model/",
			"Rebecca Bilbro",
			"Github Pages",
		},
		Story{
			"A System Health Check in Go",
			"https://rebeccabilbro.github.io/doctor-go/",
			"Rebecca Bilbro",
			"Github Pages",
		},
		Story{
			"Words in Space",
			"https://rebeccabilbro.github.io/words-in-space/",
			"Rebecca Bilbro",
			"Github Pages",
		},
		Story{
			"Getting Started with Go",
			"https://rebeccabilbro.github.io/getting-started-with-go/",
			"Rebecca Bilbro",
			"Github Pages",
		},
		Story{
			"SPARQL Queries for Local RDF Data",
			"https://rebeccabilbro.github.io/rdflib-and-sparql/",
			"Rebecca Bilbro",
			"Github Pages",
		},
	)
}

func searchStories(query string) []Story {
	var foundStories []Story
	for _, story := range stories {
		if strings.Contains(strings.ToUpper(story.title), strings.ToUpper(query)) {
			foundStories = append(foundStories, story)
		}
	}
	return foundStories
}

func search(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("q")
	if query == "" {
		http.Error(w, "Search parameter q is required to search", http.StatusNotAcceptable)
		return
	}

	w.Write([]byte("<html><body>"))
	s := searchStories(query)
	if len(s) == 0 {
		w.Write([]byte(fmt.Sprintf("No results for query '%s' .\n<br>", query)))
	} else {
		for _, story := range s {
			w.Write([]byte(fmt.Sprintf("<a href='%s'>%s</a><br>by %s on %s<br><br>", story.url, story.title, story.author, story.source)))
		}
	}
	w.Write([]byte("<a href='../'>Back</a>"))
	w.Write([]byte("</body></html>"))
}

func topTen(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<html><body>"))
	form := "<form action='search' method='get'>Search: <input type='test' name='q'> <input type='submit'></form>"
	w.Write([]byte(form))
	for i := len(stories) - 1; i >= 0 && len(stories)-i < 10; i-- {
		story := stories[i]
		w.Write([]byte(fmt.Sprintf("<a href='%s'>%s</a><br>by %s on %s<br><br>", story.url, story.title, story.author, story.source)))
	}
	w.Write([]byte("</body></html>"))
}

func main() {
	http.HandleFunc("/", topTen)
	http.HandleFunc("/search", search)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
