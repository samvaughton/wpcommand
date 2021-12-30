package registry

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	appConfig "github.com/samvaughton/wpcommand/v2/pkg/config"
	"github.com/samvaughton/wpcommand/v2/pkg/notify"
	"github.com/samvaughton/wpcommand/v2/pkg/pipeline"
	"github.com/samvaughton/wpcommand/v2/pkg/preview"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
	"text/template"
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

					bpr := types.BuildRequest{
						Id: bprc.BuildId,

						RepoOwner: appConfig.Config.Github.Owner,
						RepoName:  appConfig.Config.Github.Repo,
						Workflow:  appConfig.Config.Github.PreviewWorkflow,

						RepoRef: bprc.BuildPreviewRef,

						DockerRegistryName: appConfig.Config.Docker.PreviewRepo,
						WordpressDomain:    job.Site.WpDomain,

						IsPreviewBuild: true,
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

					if valid == false {
						return nil, errors.New("could not locate tracker object from workflow command")
					}

					ticker := time.NewTicker(1 * time.Minute)

					// deadline
					deadline := time.Now().Add(120 * time.Minute)

					for {
						select {
						case t := <-ticker.C:
							if time.Now().After(deadline) {
								return nil, errors.New("timeout on docker image check")
							}

							log.Debug(fmt.Sprintf("checking dockerhub registry for tag %s ", buildId), t.String())

							tag, err := preview.GetDockerTag(appConfig.Config.Docker.GetPreviewImageName(), buildId)

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
					previewImageName := appConfig.Config.Docker.GetPreviewImageName()

					log.Debug(fmt.Sprintf("deploying preview build %s:%s to kubernetes", previewImageName, buildId))

					err := tracker.DeployPreviewBuild(preview.TemplateContext{
						ImageName: previewImageName,
						BuildId:   buildId,
						SiteId:    pipeline.Site.Uuid,
					})

					if err != nil {
						return nil, err
					}

					return &types.CommandResult{Command: CmdPreviewBuildKubernetesDeploy, Output: buildId}, nil
				},
			},
			&pipeline.WrappedCommand{
				Name: CmdPreviewBuildNotifyAccountUsers,
				Wrapped: func(pipeline *pipeline.SiteCommandPipeline) (*types.CommandResult, error) {

					buildId := pipeline.PreviousResult.Output
					url := preview.GetPreviewUrl(buildId)

					t, err := template.New("build-preview-email").
						Parse("Hi there,\n\nA preview build is now ready for your site {{ .SiteName }}. You may find it here:\n\n{{ .Url }}\n\nThis preview build will be automatically deleted after 4 hours.\n\nKind Regards, Rentivo")

					if err != nil {
						return nil, err
					}

					var output bytes.Buffer

					err = t.Execute(&output, map[string]string{
						"Url":      url,
						"BuildId":  buildId,
						"SiteName": pipeline.Site.Description,
					})

					if err != nil {
						return nil, err
					}

					err = notify.SendToAllUsers(
						fmt.Sprintf("%s - Preview Build Ready", pipeline.Site.Description),
						output.String(),
						pipeline.Site.AccountId,
					)

					if err != nil {
						return nil, err
					}

					return &types.CommandResult{Command: CmdPreviewBuildNotifyAccountUsers, Output: url}, nil
				},
			},
		},
	}
}
