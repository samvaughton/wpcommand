package server

import (
	"embed"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/samvaughton/wpcommand/v2/pkg/auth"
	"github.com/samvaughton/wpcommand/v2/pkg/config"
	"github.com/samvaughton/wpcommand/v2/pkg/db"
	"github.com/samvaughton/wpcommand/v2/pkg/flow"
	"github.com/samvaughton/wpcommand/v2/pkg/registry"
	"github.com/samvaughton/wpcommand/v2/pkg/scheduler"
	"github.com/samvaughton/wpcommand/v2/pkg/util"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func Start(staticFiles *embed.FS, configData string, authData string) {
	config.InitConfig(configData)
	db.InitDbConnection()
	util.SetupLogging()

	auth.InitAuth(authData)

	flow.CreateDefaultAccountAndUser()

	registry.CreateDefaultCommands()

	scheduler.Init(time.Second*5, 2)
	scheduler.Start()
	scheduler.SetupCron()

	router := mux.NewRouter()

	SetupApi(router)
	SetupPublic(router)
	SetupSpa(router, staticFiles)

	srv := &http.Server{
		Handler:      router,
		Addr:         config.Config.ServerAddress,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Info(fmt.Sprintf("api server listening on %s", srv.Addr))

	log.Fatal(srv.ListenAndServe())
}
