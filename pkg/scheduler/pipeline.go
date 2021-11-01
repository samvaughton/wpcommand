package scheduler

import (
	"context"
	"fmt"
	"github.com/samvaughton/wpcommand/v2/pkg/config"
	"github.com/samvaughton/wpcommand/v2/pkg/db"
	"github.com/samvaughton/wpcommand/v2/pkg/execution"
	"github.com/samvaughton/wpcommand/v2/pkg/pipeline"
	"github.com/samvaughton/wpcommand/v2/pkg/registry"
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

					cmdConfig := map[string]interface{}{}

					if local.Command.Type == types.CommandTypePreviewBuild {
						// load up a request job from the job log
					}

					commands := GetCommandsFromJob(local, cmdConfig)

					AddPipelineTask(pipeline.SiteCommandPipeline{
						Site:     local.Site,
						Executor: executor,
						Config:   config.Config,
						Hooks: pipeline.CommandHooks{
							Started: []func(){
								func() {
									meta := make(map[string]interface{})

									local.Status = types.CommandJobStatusRunning

									db.CreateCommandJobEvent(
										local.Id,
										types.EventLogTypeJobStarted,
										types.EventLogStatusSuccess,
										"job",
										fmt.Sprintf("root_level_commands=%v", len(commands)),
										meta,
									)

									_, err = db.Db.NewUpdate().Model(&local).WherePK().Exec(context.Background())
								},
							},
							Finished: []func([]error){
								func(errors []error) {
									meta := make(map[string]interface{})

									if len(errors) > 0 {
										local.Status = types.CommandJobStatusFailure
									} else {
										local.Status = types.CommandJobStatusSuccess
									}

									db.CreateCommandJobEvent(
										local.Id,
										types.EventLogTypeJobFinished,
										types.EventLogStatusSuccess,
										"job",
										fmt.Sprintf("errors=%v", len(errors)),
										meta,
									)

									_, err = db.Db.NewUpdate().Model(&local).WherePK().Exec(context.Background())
								},
							},
							PostSuccess: []func(pipeline.SiteCommand, *types.CommandResult, error){
								func(c pipeline.SiteCommand, result *types.CommandResult, err error) {
									meta := make(map[string]interface{})

									jType := types.EventLogTypeInfo
									output := ""
									if result != nil {
										//output = result.Output

										if result.Data != nil {
											jType = types.EventLogTypeData
											meta["Data"] = result.Data
										}
									}

									_, err = db.CreateCommandJobEvent(
										local.Id,
										jType,
										types.EventLogStatusSuccess,
										fmt.Sprintf("job.%s", c.GetName()),
										output,
										meta,
									)

									if err != nil {
										log.Error(err)
									}
								},
							},
							// run when a pre-check command fails
							Skipped: []func(pipeline.SiteCommand, *types.CommandResult, error){
								func(c pipeline.SiteCommand, result *types.CommandResult, err error) {
									meta := make(map[string]interface{})

									jType := types.EventLogTypeInfo
									output := ""
									if result != nil {
										output = result.Output

										if result.Data != nil {
											jType = types.EventLogTypeData
											meta["Data"] = result.Data
										}
									}

									db.CreateCommandJobEvent(
										local.Id,
										jType,
										types.EventLogStatusSkipped,
										fmt.Sprintf("job.%s", c.GetName()),
										output,
										meta,
									)
								},
							},
							PostError: []func(pipeline.SiteCommand, *types.CommandResult, error){
								func(c pipeline.SiteCommand, result *types.CommandResult, err error) {
									meta := make(map[string]interface{})

									jType := types.EventLogTypeInfo
									if result != nil {
										if result.Data != nil {
											jType = types.EventLogTypeData
											meta["Data"] = result.Data
										}
									}

									db.CreateCommandJobEvent(
										local.Id,
										jType,
										types.EventLogStatusFailure,
										fmt.Sprintf("job.%s", c.GetName()),
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

func GetCommandsFromJob(job types.CommandJob, cmdConfig map[string]interface{}) []pipeline.SiteCommand {
	commands := make([]pipeline.SiteCommand, 0)

	switch job.Command.Type {
	case types.CommandTypeWpBuiltIn:
		commands = []pipeline.SiteCommand{
			registry.BuiltInCommandRegistry[job.Key](job.Site, cmdConfig),
		}
		break
	case types.CommandTypeHttpCall:
		commands = []pipeline.SiteCommand{
			registry.GetHttpCallCommand(job),
		}
	case types.CommandTypePreviewBuild:
		commands = []pipeline.SiteCommand{
			registry.GetPreviewBuildCommand(job),
		}
	case types.CommandTypeBuildRelease:
		commands = []pipeline.SiteCommand{
			registry.GetBuildReleaseCommand(job),
		}
	}

	return commands
}
