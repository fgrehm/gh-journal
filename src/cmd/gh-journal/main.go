package main

import (
	"os"

	"ghjournal"

	log "github.com/Sirupsen/logrus"
	"github.com/jasonlvhit/gocron"
)

func main() {
	go ghjournal.SyncEvents()
	gocron.Every(5).Minutes().Do(ghjournal.SyncEvents)
	gocron.Start()

	log.SetLevel(log.InfoLevel)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	ghjournal.RunServer(port)
}
