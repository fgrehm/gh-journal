package main

import (
	"fmt"
	"os"
	"time"

	"ghjournal"

	"github.com/jasonlvhit/gocron"
	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

func main() {
	syncAndShowStars()
	gocron.Every(5).Minutes().Do(syncAndShowStars)
	gocron.Start()

	log.SetLevel(log.InfoLevel)
	ghjournal.RunServer("8080")
}

func syncAndShowStars() {
	// MONGO_HOST
	mongoHost := os.Getenv("MONGO_PORT_27017_TCP_ADDR")
	token := os.Getenv("GITHUB_TOKEN")
	userName := "fgrehm"

	client, err := ghjournal.NewGitHubClient(userName, token)
	if err != nil {
		panic(err)
	}

	session, err := mgo.Dial(mongoHost)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	repo := ghjournal.NewEventsRepository(session.DB("gh-journal"))
	if err = ghjournal.ImportEvents(client, repo); err != nil {
		panic(err)
	}

	beginningOfDay, _ := time.Parse(time.RFC3339, "2015-11-28T00:00:00Z")
	endOfDay, _ := time.Parse(time.RFC3339, "2015-12-31T23:59:59Z")
	log.Printf("Building `%s` edition", beginningOfDay.Format("2006-01-02"))
	events, err := repo.EventsWithin(beginningOfDay, endOfDay)
	if err != nil {
		panic(err)
	}
	for _, event := range events {
		eventAction := "<nil>"
		if event.Action != nil {
			eventAction = *event.Action
		}
		project := fmt.Sprintf("%s/%s", event.Project.Owner, event.Project.Name)

		if event.ProjectStarred() {
			log.Printf(" * '%s' was starred by '%s'", project, event.Actor)
		} else {
			log.Debugf("TODO: Handle Type=%s, Action=%s, Project=%s, Actor=%s", event.Type, eventAction, project, event.Actor)
		}
	}
}
