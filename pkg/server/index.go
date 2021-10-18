package server

import (
	"context"
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
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Start(staticFiles *embed.FS, configData string, authData string) {
	config.InitConfig(configData)
	db.InitDbConnection()
	util.SetupLogging()

	auth.InitAuth(authData)

	flow.CreateDefaultAccountAndUser()

	registry.CreateDefaultCommands()

	// nearly all tasks will be http blocking hence a high number of workers
	scheduler.Init(time.Second*5, 50)
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

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Info(fmt.Sprintf("api server listening on %s", srv.Addr))

	<-done
	log.Info("server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here

		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %+v", err)
	}

	log.Info("server exited properly")
}
