package main

import (
	"ghjournal"
	"gopkg.in/mgo.v2"
	"log"
	"os"
)

func main() {
	token := os.Getenv("GITHUB_TOKEN")
	userName := "fgrehm"
	client := ghjournal.NewGitHubClient(userName, token)

	session, err := mgo.Dial("172.17.0.2")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	eventsCollection := session.DB("gh-journal").C("events")
	for pageNum := 1; pageNum <= 10; pageNum++ {
		log.Printf("Fetching events from page %d...", pageNum)
		ghEvents, err := client.Events(pageNum)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%+v", len(ghEvents))
		for _, event := range ghEvents {
			err = eventsCollection.Insert(event)
			if err != nil {
				log.Fatal(err)
			}
		}
		log.Printf("Done fetching events from page %d", pageNum)
	}
}
