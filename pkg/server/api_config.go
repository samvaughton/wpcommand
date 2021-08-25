package server

import (
	"encoding/json"
	"github.com/samvaughton/wpcommand/v2/pkg/config"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func configHandler(resp http.ResponseWriter, req *http.Request) {
	bytes, err := json.Marshal(config.Config)

	if err != nil {
		log.WithFields(log.Fields{"endpoint": "/api/config"}).Error(err)

		resp.WriteHeader(500)

		return
	}

	_, err = resp.Write(bytes)

	if err != nil {
		log.WithFields(log.Fields{"endpoint": "/api/config"}).Error(err)
	}
}
