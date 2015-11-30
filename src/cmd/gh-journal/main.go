package main

import (
	"os"

	"ghjournal"
	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

func main() {
	log.SetLevel(log.DebugLevel)

	token := os.Getenv("GITHUB_TOKEN")
	userName := "fgrehm"

	client, err := ghjournal.NewGitHubClient(userName, token)
	if err != nil {
		panic(err)
	}

	session, err := mgo.Dial("172.17.0.2")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	repo := ghjournal.NewEventsRepository(session.DB("gh-journal"))
	if err = ghjournal.ImportEvents(client, repo); err != nil {
		panic(err)
	}
}
