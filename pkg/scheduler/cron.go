package scheduler

import (
	"github.com/robfig/cron"
	"github.com/samvaughton/wpcommand/v2/pkg/db"
	"github.com/samvaughton/wpcommand/v2/pkg/registry"
	log "github.com/sirupsen/logrus"
)

var Cron *cron.Cron

func SetupCron() {
	Cron = cron.New()

	_ = Cron.AddFunc("* */30 * * * *", func() {
		log.WithFields(log.Fields{
			"Source": "CRON",
			"Action": "WP_USER_SYNC",
			"Detail": "",
		}).Info("Cron starting")

		// locate the user sync command
		command, err := db.CommandGetByKey(registry.CmdWpSiteUserSync)

		if err != nil {
			log.WithFields(log.Fields{
				"Source": "CRON",
				"Action": "WP_USER_SYNC",
				"Detail": "GET_COMMAND",
			}).Error(err)

			return
		}

		// we need to create a command job for each site to sync the users

		sites, err := db.SiteGetAllEnabled()

		if err != nil {
			log.WithFields(log.Fields{
				"Source": "CRON",
				"Action": "WP_USER_SYNC",
				"Detail": "LIST_SITES",
			}).Error(err)

			return
		}

		db.CreateCommandJobs(command, sites, 0) // system
	})

	Cron.Start()
	log.Info("cron initialized")
}
