package ghjournal

import (
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

type Event struct {
  ID string
  Type string
  Data map[string]interface{}
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

  event := Event {
    ID: id,
    Type: eventType,
    Data: ghEvent,
  }
  return r.collection.Insert(event)
}
