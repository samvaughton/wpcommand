package scheduler

import (
	"github.com/robfig/cron/v3"
	"github.com/samvaughton/wpcommand/v2/pkg/flow"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
)

var Cron *cron.Cron

func SetupCron() {
	Cron = cron.New()

	Cron.AddFunc("*/15 * * * *", func() {
		flow.RunWpUserSync(types.FlowOptions{LogSource: "CRON"})
	})

	Cron.Start()

	log.Info("cron initialized")
}
