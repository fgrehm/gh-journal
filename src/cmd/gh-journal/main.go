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

	repo := ghjournal.NewEventsRepository(session.DB("gh-journal"))
	for pageNum := 1; pageNum <= 10; pageNum++ {
		log.Printf("Importing events from page %d...", pageNum)
		ghEvents, err := client.Events(pageNum)
		if err != nil {
			log.Fatal(err)
		}
		createdAtLeastOneEvent := false
		for _, event := range ghEvents {
			id := event["id"].(string)
			if exists, err := repo.Exists(id); err != nil {
				log.Fatal(err)
			} else if exists {
				log.Printf("SKIP ID=%s", id)
				continue
			}

			createdAtLeastOneEvent = true
			log.Printf("CREATE ID=%s", id)
			err = repo.Insert(event)
			if err != nil {
				log.Fatal(err)
			}
		}
		if !createdAtLeastOneEvent {
			log.Printf("Done with import at page %d", pageNum)
			break
		}
	}
}
