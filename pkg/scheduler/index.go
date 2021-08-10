package scheduler

import (
	"context"
	"fmt"
	"github.com/samvaughton/wpcommand/v2/pkg/config"
	"github.com/samvaughton/wpcommand/v2/pkg/db"
	"github.com/samvaughton/wpcommand/v2/pkg/execution"
	"github.com/samvaughton/wpcommand/v2/pkg/pipeline"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	"github.com/samvaughton/wpcommand/v2/pkg/worker"
	log "github.com/sirupsen/logrus"
	"time"
)

var done chan bool
var ticker *time.Ticker
var workerPool *worker.CommandRunnerPool

func Init(interval time.Duration, maxWorkers int) {
	ticker = time.NewTicker(interval)
	done = make(chan bool)
	workerPool = worker.NewCommandRunnerPool(maxWorkers)
}

func AddPipelineTask(scp pipeline.SiteCommandPipeline) {
	workerPool.AddTask(scp)
}

func Start() {
	workerPool.Start()

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				// move created jobs to the pending status and send them to the command runners
				log.Debug("checking pending jobs")

				created := db.GetCreatedJobs()

				if len(created) > 0 {
					log.Infof("found %v jobs", len(created))
				}

				for _, item := range created {
					local := item
					executor, err := execution.NewCommandExecutor(local.Site)

					if err != nil {
						// set the status to failed, job to failed and add an event describing situation
						local.Status = types.CommandJobStatusFailure
						_, err = db.Db.NewUpdate().Model(&local).WherePK().Exec(context.Background())

						log.Error(err)
						continue
					}

					local.Status = types.CommandJobStatusPending
					_, err = db.Db.NewUpdate().Model(&local).WherePK().Exec(context.Background())

					if err != nil {
						log.Error(err)

						continue
					}

					log.Infof("dispatching job %s to pool", item.Uuid)

					commands := []pipeline.SiteCommand{
						pipeline.CommandRegistry[local.Key](local.Site),
					}

					AddPipelineTask(pipeline.SiteCommandPipeline{
						Site:     local.Site,
						Executor: executor,
						Config:   config.Config,
						Hooks: pipeline.CommandHooks{
							Started: []func(){
								func() {
									meta := map[string]interface{}{}

									local.Status = types.CommandJobStatusRunning

									db.CreateCommandJobEvent(
										local.Id,
										"",
										types.EventLogStatusJobStarted,
										"",
										fmt.Sprintf("root_level_commands=%v", len(commands)),
										meta,
									)

									_, err = db.Db.NewUpdate().Model(&local).WherePK().Exec(context.Background())
								},
							},
							Finished: []func([]error){
								func(errors []error) {
									meta := map[string]interface{}{}

									if len(errors) > 0 {
										local.Status = types.CommandJobStatusFailure
									} else {
										local.Status = types.CommandJobStatusSuccess
									}

									db.CreateCommandJobEvent(
										local.Id,
										"",
										types.EventLogStatusJobFinished,
										"",
										fmt.Sprintf("errors=%v", len(errors)),
										meta,
									)

									_, err = db.Db.NewUpdate().Model(&local).WherePK().Exec(context.Background())
								},
							},
							PostSuccess: []func(pipeline.SiteCommand, *types.CommandResult, error){
								func(c pipeline.SiteCommand, result *types.CommandResult, err error) {
									meta := map[string]interface{}{}

									command := ""
									output := ""

									if result != nil {
										command = result.Command
									}

									db.CreateCommandJobEvent(
										local.Id,
										c.GetName(),
										types.EventLogStatusSuccess,
										command,
										output,
										meta,
									)
								},
							},
							// run when a pre-check command fails
							Skipped: []func(pipeline.SiteCommand, *types.CommandResult, error){
								func(c pipeline.SiteCommand, result *types.CommandResult, err error) {
									meta := map[string]interface{}{}

									command := ""

									if result != nil {
										command = result.Command
									}

									db.CreateCommandJobEvent(
										local.Id,
										c.GetName(),
										types.EventLogStatusSkipped,
										command,
										fmt.Sprintf("%s", err),
										meta,
									)
								},
							},
							PostError: []func(pipeline.SiteCommand, *types.CommandResult, error){
								func(c pipeline.SiteCommand, result *types.CommandResult, err error) {
									meta := map[string]interface{}{}

									command := ""

									if result != nil {
										command = result.Command
									}

									db.CreateCommandJobEvent(
										local.Id,
										c.GetName(),
										types.EventLogStatusFailure,
										command,
										fmt.Sprintf("%s", err),
										meta,
									)
								},
							},
						},
						Options:  pipeline.ExecuteOptions{},
						Commands: commands,
					})
				}
			}
		}
	}()
}
