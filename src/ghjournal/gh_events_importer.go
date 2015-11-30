package ghjournal

import (
  log "github.com/Sirupsen/logrus"
)

func ImportEvents(client GitHubClient, repo EventsRepository) error {
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
