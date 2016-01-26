package ghjournal

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

func SyncEvents() {
	// MONGO_HOST
	mongoHost := os.Getenv("MONGO_PORT_27017_TCP_ADDR")
	token := os.Getenv("GITHUB_TOKEN")
	userName := os.Getenv("GITHUB_USER")
	if userName == "" {
		panic("GITHUB_USER is not set")
	}

	client, err := NewGitHubClient(userName, token)
	if err != nil {
		panic(err)
	}

	session, err := mgo.Dial(mongoHost)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	repo := NewEventsRepository(session.DB("gh-journal"))
	if err = importEvents(client, repo); err != nil {
		panic(err)
	}
}

func importEvents(client GitHubClient, repo EventsRepository) error {
	for pageNum := 1; pageNum <= 10; pageNum++ {
		log.Printf("Importing events from page %d...", pageNum)
		ghEvents, err := client.Events(pageNum)
		if err != nil {
			return err
		}
		createdAtLeastOneEvent := false
		for _, event := range ghEvents {
			id := event["id"].(string)
			if exists, err := repo.Exists(id); err != nil {
				return err
			} else if exists {
				log.Debugf("SKIP ID=%s", id)
				continue
			}

			createdAtLeastOneEvent = true
			log.Printf("CREATE ID=%s", id)
			err = repo.Insert(event)
			if err != nil {
				return err
			}
		}
		if !createdAtLeastOneEvent {
			log.Printf("Done with import at page %d", pageNum)
			break
		}
	}
	return nil
}
