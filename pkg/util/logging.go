package util

import (
	"github.com/samvaughton/wpcommand/v2/pkg/config"
	log "github.com/sirupsen/logrus"
)

func SetupLogging() {
	if config.Config.Environment == "" {
		log.Fatalf("environment not set")
	}

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	level, err := log.ParseLevel(config.Config.Logging.Level)

	if err != nil {
		log.Fatalf("could not parse logging level")
	}

	log.SetLevel(level)

	log.WithFields(log.Fields{
		"event": "INITIALIZATION",
	}).Debug("initialized")
}
