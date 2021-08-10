package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/samvaughton/wpcommand/v2/pkg/config"
	"github.com/samvaughton/wpcommand/v2/pkg/db"
	"github.com/samvaughton/wpcommand/v2/pkg/logutil"
	"github.com/samvaughton/wpcommand/v2/pkg/scheduler"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func Start() {
	config.InitConfig()
	db.InitDbConnection()
	logutil.SetupLogging()

	scheduler.Init(time.Second*5, 1) // check every 5 seconds
	scheduler.Start()

	router := mux.NewRouter()

	SetupApi(router)
	SetupSpa(router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":8999",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Info(fmt.Sprintf("api server listening on %s", srv.Addr))

	log.Fatal(srv.ListenAndServe())
}
