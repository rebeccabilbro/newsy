package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/caser/gophernews"
	"github.com/jzelinskie/geddit"
)

var redditSession *geddit.LoginSession
var hackerNewsClient *gophernews.Client
var stories []Story

func init() {
	hackerNewsClient = gophernews.NewClient()
	var err error
    redditSession, err = geddit.NewLoginSession("redditUsername", "redditPassword", "customUserAgent")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type Story struct {
	title  string
	url    string
	author string
	source string
}

func getHnStoryDetails(id int, c chan<- Story, wg *sync.WaitGroup) {
	defer wg.Done()
	story, err := hackerNewsClient.GetStory(id)
	if err != nil {
		return
	}
	newStory := Story{
		title:  story.Title,
		url:    story.URL,
		author: story.By,
		source: "HackerNews",
	}
	c <- newStory
}

func newHnStories(c chan<- Story) {
	defer close(c)
	changes, err := hackerNewsClient.GetChanges()
	if err != nil {
		fmt.Println(err)
		return
	}
	var wg sync.WaitGroup
	for _, id := range changes.Items {
		wg.Add(1)
		go getHnStoryDetails(id, c, &wg)
	}
	wg.Wait()
}

func newRedditStories(c chan<- Story) {
	defer close(c)
	sort := geddit.PopularitySort(geddit.NewSubmissions)
	var listingOptions geddit.ListingOptions
	submissions, err := redditSession.SubredditSubmissions("programming", sort, listingOptions)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, s := range submissions {
		newStory := Story{
			title:  s.Title,
			url:    s.URL,
			author: s.Author,
			source: "Reddit /r/programming",
		}

		c <- newStory
	}
}

func outputToStories(c <-chan Story) {
	for s := range c {
		stories = append(stories, s)
	}
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
	go func() {
		for {
			fmt.Println("Fetching the latest news...")
			fromHn := make(chan Story, 8)
			fromReddit := make(chan Story, 8)
			toList := make(chan Story, 8)
			go outputToStories(toList)
			go newHnStories(fromHn)
			go newRedditStories(fromReddit)

			hnOpen := true
			redditOpen := true

			for hnOpen || redditOpen {
				select {
				case story, open := <-fromHn:
					if open {
						toList <- story
					} else {
						hnOpen = false
					}
				case story, open := <-fromReddit:
					if open {
						toList <- story
					} else {
						redditOpen = false
					}
				}
			}
			fmt.Println("Got all the new stories!")
			time.Sleep(30 * time.Second)
		}
	}()

	http.HandleFunc("/", topTen)
	http.HandleFunc("/search", search)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
