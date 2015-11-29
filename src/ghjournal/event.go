package ghjournal

import (
  "time"
)

type Event struct {
  ID string
  Type string
  Action *string
  CreatedAt time.Time `bson:"created_at"`
  Project Project
  Raw map[string]interface{}
}

type Project struct {
  Owner string
  Name  string
}
