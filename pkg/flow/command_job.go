package flow

import (
	"context"
	"fmt"
	"github.com/samvaughton/wpcommand/v2/pkg/db"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
	"time"
)

func CleanupAbandonedJobs(flowOpts types.FlowOptions) {
	log.WithFields(log.Fields{
		"Source": flowOpts.LogSource,
		"Action": "CLEANUP_ABANDONED_JOBS",
		"Detail": "START",
	}).Debug("started")

	jobs, err := db.CommandJobsGetAbandoned(1 * time.Hour)

	if err != nil {
		log.WithFields(log.Fields{
			"Source": flowOpts.LogSource,
			"Action": "CLEANUP_ABANDONED_JOBS",
			"Detail": "GET_JOBS",
		}).Debug(err)

		return
	}

	log.WithFields(log.Fields{
		"Source": flowOpts.LogSource,
		"Action": "CLEANUP_ABANDONED_JOBS",
		"Detail": "GET_JOBS",
	}).Debug(fmt.Sprintf("found %v abandoned jobs", len(jobs)))

	for _, job := range jobs {
		job.Status = types.CommandJobStatusTerminated
		db.Db.NewUpdate().Model(job).WherePK().Returning("*").Exec(context.Background())

		_, err := db.CreateCommandJobEvent(job.Id, types.EventLogTypeInfo, types.EventLogStatusTerminated, "job", "job terminated due to timeout via cron", map[string]interface{}{})

		if err != nil {
			log.WithFields(log.Fields{
				"Source": flowOpts.LogSource,
				"Action": "CLEANUP_ABANDONED_JOBS",
				"Detail": "CREATE_COMMAND_JOB_EVENT",
			}).Error(err)

			continue
		}

		log.WithFields(log.Fields{
			"Source": flowOpts.LogSource,
			"Action": "CLEANUP_ABANDONED_JOBS",
			"Detail": "JOB_CLEANUP",
			"JobId":  job.Uuid,
		}).Info("cleaned up abandoned job")
	}

	log.WithFields(log.Fields{
		"Source": flowOpts.LogSource,
		"Action": "CLEANUP_ABANDONED_JOBS",
		"Detail": "FINISH",
	}).Debug("finished")
}
