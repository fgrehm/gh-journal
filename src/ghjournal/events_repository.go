package ghjournal

import (
  "strings"
  "time"

  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

type Project struct {
  Owner string
  Name  string
}

type Event struct {
  ID string
  Type string
  CreatedAt time.Time `bson:"created_at"`
  Project Project
  Raw map[string]interface{}
}

type EventsRepository interface {
  Exists(id string) (bool, error)
  Insert(GitHubEvent) error
}

type eventsRepository struct {
  collection *mgo.Collection
}

func NewEventsRepository(db *mgo.Database) EventsRepository {
  return &eventsRepository{db.C("events")}
}

func (r *eventsRepository) Exists(id string) (bool, error) {
  query := r.collection.Find(bson.M{"id": id})
  count, err := query.Count()
  return count > 0, err
}

func (r *eventsRepository) Insert(ghEvent GitHubEvent) error {
  id := ghEvent["id"].(string)
  delete(ghEvent, "id")

  eventType := ghEvent["type"].(string)
  delete(ghEvent, "type")

  createdAt, err := time.Parse(time.RFC3339, ghEvent["created_at"].(string))
  if err != nil {
    return err
  }
  delete(ghEvent, "created_at")

  repoSlug := ghEvent["repo"].(map[string]interface{})["name"].(string)
  projectOwnerAndName := strings.SplitN(repoSlug, "/", 2)
  project := Project{
    Owner: projectOwnerAndName[0],
    Name:  projectOwnerAndName[1],
  }

  event := Event {
    ID: id,
    Type: eventType,
    CreatedAt: createdAt,
    Project: project,
    Raw: ghEvent,
  }
  return r.collection.Insert(event)
}
