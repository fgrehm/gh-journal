package ghjournal

import (
  "strings"
  "time"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

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
  e, err := r.buildEvent(ghEvent)
  if err != nil {
    return err
  }
  return r.collection.Insert(e)
}

func (r *eventsRepository) buildEvent(ghEvent GitHubEvent) (*Event, error) {
  event := &Event {
    ID: ghEvent["id"].(string),
    Type: ghEvent["type"].(string),
    Raw: ghEvent,
  }

  createdAt, err := time.Parse(time.RFC3339, ghEvent["created_at"].(string))
  if err != nil {
    return &Event{}, err
  }
  event.CreatedAt = createdAt

  repoSlug := ghEvent["repo"].(map[string]interface{})["name"].(string)
  projectOwnerAndName := strings.SplitN(repoSlug, "/", 2)
  event.Project = Project{
    Owner: projectOwnerAndName[0],
    Name:  projectOwnerAndName[1],
  }

  var action *string = nil
  payload := ghEvent["payload"].(map[string]interface{})
  if payload["action"] != nil {
    str := payload["action"].(string)
    action = &str
  }
  event.Action = action

  return event, nil
}
