package main

import (
    "os"
    "fmt"

    "github.com/caser/gophernews"
    "github.com/jzelinskie/geddit"
)

var redditSession *geddit.LoginSession
var hackerNewsClient *gophernews.Client

func init() {
    hackerNewsClient = gophernews.NewClient()

    // todo: switch to OAuth2 method
    // see https://github.com/jzelinskie/geddit/blob/master/example_test.go
    var err error
    redditSession, err = geddit.NewLoginSession("redditUsername", "redditPassword", "customUserAgent")

    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

type Story struct {
    title string
    url string
    author string
    source string
}

func newHnStories() []Story {
    var stories []Story
    changes, err := hackerNewsClient.GetChanges()
    if err != nil {
        fmt.Println(err)
        return nil
    }

    for _, id := range changes.Items {
        story, err := hackerNewsClient.GetStory(id)
        if err != nil {
            continue
        }

        newStory := Story{
            title: story.Title,
            url: story.URL,
            author: story.By,
            source: "HackerNews",
        }

        stories = append(stories, newStory)
    }

    return stories
}

func newRedditStories() []Story {
    var stories []Story
    sort := geddit.PopularitySort(geddit.NewSubmissions)
    var listingOptions geddit.ListingOptions
    submissions, err := redditSession.SubredditSubmissions("programming", sort, listingOptions)
    if err != nil {
        fmt.Println(err)
        return nil
    }

    for _, s := range submissions {
        newStory := Story {
            title: s.Title,
            url: s.URL,
            author: s.Author,
            source: "Reddit /r/programming",
        }

        stories = append(stories, newStory)
    }

    return stories
}

func main() {
    hnStories :=newHnStories()
    redditStories := newRedditStories()

    var stories []Story

    if hnStories != nil {
        stories = append(stories, hnStories...)
    }

    if redditStories != nil {
        stories = append(stories, redditStories...)
    }

    file, err := os.Create("data/stories.txt")
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer file.Close()

    for _, s := range stories {
        fmt.Fprintf(file, "%s: %s \nby %s on %s\n\n", s.title, s.url, s.author, s.source)
    }
}
