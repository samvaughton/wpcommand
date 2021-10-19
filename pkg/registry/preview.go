package registry

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	appConfig "github.com/samvaughton/wpcommand/v2/pkg/config"
	"github.com/samvaughton/wpcommand/v2/pkg/pipeline"
	"github.com/samvaughton/wpcommand/v2/pkg/preview"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
	"time"
)

func GetPreviewBuildCommand(job types.CommandJob) pipeline.SiteCommand {
	return &pipeline.SimplePipelineCommand{
		Name: CmdPreviewBuild,
		Commands: []pipeline.SiteCommand{
			&pipeline.WrappedCommand{
				Name: CmdPreviewBuildRunGithubWorkflow,
				Wrapped: func(pipeline *pipeline.SiteCommandPipeline) (*types.CommandResult, error) {
					bprc, err := types.NewCommandTypePreviewBuildConfigFromConfig(job.Config)

					if err != nil {
						return nil, err
					}

					bpr := types.BuildPreviewRequest{
						Id: uuid.New().String(),

						RepoOwner: appConfig.Config.Github.Owner,
						RepoName:  appConfig.Config.Github.Repo,
						Workflow:  appConfig.Config.Github.PreviewWorkflow,

						RepoRef: bprc.BuildPreviewRef,

						DockerRegistryName: job.Site.DockerRegistryName,
						WordpressDomain:    job.Site.WpDomain,
					}

					tracker := preview.NewBuildTracker(preview.NewGithubClient(appConfig.Config.GithubToken), appConfig.Config.K8RestConfig)

					err = tracker.RunGithubWorkflowJob(bpr)

					if err != nil {
						return nil, err
					}

					return &types.CommandResult{Command: CmdPreviewBuildRunGithubWorkflow, Output: bpr.Id, Data: tracker}, nil
				},
			},
			&pipeline.WrappedCommand{
				Name: CmdPreviewBuildRunDockerRegistryWorkflow,
				Wrapped: func(pipeline *pipeline.SiteCommandPipeline) (*types.CommandResult, error) {
					// We need to wait on successful dockerhub upload before we can deploy to K8
					tracker, valid := pipeline.PreviousResult.Data.(*preview.BuildTracker)

					buildId := pipeline.PreviousResult.Output

					repo := fmt.Sprintf(
						"%s/%s",
						appConfig.Config.Docker.Namespace,
						job.Site.DockerRegistryName,
					)

					if valid == false {
						return nil, errors.New("could not locate tracker object from workflow command")
					}

					// check every 15 seconds against the registry
					ticker := time.NewTicker(30 * time.Second)

					// deadline
					deadline := time.Now().Add(20 * time.Minute)

					for {
						select {
						case t := <-ticker.C:
							if time.Now().After(deadline) {
								return nil, errors.New("timeout on docker image check")
							}

							log.Debug(fmt.Sprintf("checking dockerhub registry for tag %s ", buildId), t.String())

							tag, err := preview.GetDockerTag(repo, buildId)

							if err != nil {
								log.Error(err)
							}

							if tag != nil {
								return &types.CommandResult{Command: CmdPreviewBuildRunDockerRegistryWorkflow, Output: buildId, Data: tracker}, nil
							}
						}
					}
				},
			},
			&pipeline.WrappedCommand{
				Name: CmdPreviewBuildKubernetesDeploy,
				Wrapped: func(pipeline *pipeline.SiteCommandPipeline) (*types.CommandResult, error) {
					// We need to wait on successful dockerhub upload before we can deploy to K8
					tracker, valid := pipeline.PreviousResult.Data.(*preview.BuildTracker)

					if valid == false {
						return nil, errors.New("could not locate tracker object from workflow command")
					}

					buildId := pipeline.PreviousResult.Output

					// create the actual output required for K8 deployment (the image)
					imageName := fmt.Sprintf("%s/%s:%s", appConfig.Config.Docker.Namespace, job.Site.DockerRegistryName, buildId)

					log.Debug(fmt.Sprintf("deploying preview build %s to kubernetes", imageName))

					err := tracker.DeployPreviewBuild(imageName, buildId)

					if err != nil {
						return nil, err
					}

					return &types.CommandResult{Command: CmdPreviewBuildKubernetesDeploy, Output: imageName, Data: tracker}, nil
				},
			},
		},
	}
}
