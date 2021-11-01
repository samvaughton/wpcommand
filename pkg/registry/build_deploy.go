package registry

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/google/uuid"
	appConfig "github.com/samvaughton/wpcommand/v2/pkg/config"
	"github.com/samvaughton/wpcommand/v2/pkg/notify"
	"github.com/samvaughton/wpcommand/v2/pkg/pipeline"
	"github.com/samvaughton/wpcommand/v2/pkg/preview"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
	"text/template"
	"time"
)

func GetBuildReleaseCommand(job types.CommandJob) pipeline.SiteCommand {
	return &pipeline.SimplePipelineCommand{
		Name: CmdBuildRelease,
		Commands: []pipeline.SiteCommand{
			&pipeline.WrappedCommand{
				Name: CmdBuildReleaseRunGithubWorkflow,
				Wrapped: func(pipeline *pipeline.SiteCommandPipeline) (*types.CommandResult, error) {
					brc, err := types.NewCommandTypeBuildReleaseConfigFromConfig(job.Config)

					if err != nil {
						return nil, err
					}

					bpr := types.BuildRequest{
						Id: uuid.New().String(),

						RepoOwner: appConfig.Config.Github.Owner,
						RepoName:  appConfig.Config.Github.Repo,
						Workflow:  brc.GithubActionName,

						RepoRef: "master",

						DockerRegistryName: job.Site.DockerRegistryName,
						WordpressDomain:    job.Site.WpDomain,
					}

					tracker := preview.NewBuildTracker(preview.NewGithubClient(appConfig.Config.GithubToken), appConfig.Config.K8RestConfig)

					err = tracker.RunGithubWorkflowJob(bpr)

					if err != nil {
						return nil, err
					}

					return &types.CommandResult{Command: CmdBuildReleaseRunGithubWorkflow, Output: bpr.Id, Data: tracker}, nil
				},
			},
			&pipeline.WrappedCommand{
				Name: CmdBuildReleaseCheckDeploymentStatus,
				Wrapped: func(pipeline *pipeline.SiteCommandPipeline) (*types.CommandResult, error) {
					// we need to check the latest on the site docker repo

					// check every 15 seconds against the registry
					ticker := time.NewTicker(30 * time.Second)

					// deadline
					deadline := time.Now().Add(20 * time.Minute)

					now := time.Now()

					for {
						select {
						case t := <-ticker.C:
							if time.Now().After(deadline) {
								return nil, errors.New("status check timed out - the build may still be running or have failed")
							}

							imageName := fmt.Sprintf("%s/%s", appConfig.Config.Docker.Namespace, pipeline.Site.DockerRegistryName)

							log.Debug("checking dockerhub registry for latest update ", t.String())

							tag, err := preview.GetDockerTag(imageName, "latest")

							if err != nil {
								log.Error(err)
							}

							// make sure the latest tag has actually been updated
							if tag != nil && tag.TagLastPushed.After(now) {
								return &types.CommandResult{Command: CmdBuildReleaseCheckDeploymentStatus, Output: tag.Name}, nil
							}
						}
					}
				},
			},
			&pipeline.WrappedCommand{
				Name: CmdBuildReleaseNotifyAccountUsers,
				Wrapped: func(pipeline *pipeline.SiteCommandPipeline) (*types.CommandResult, error) {

					buildId := pipeline.PreviousResult.Output

					t, err := template.New("build-release-email").
						Parse("Hi there,\n\nA new build has been completed for site {{ .SiteName }}. It may take a few minutes for the new build to become live.\n\nKind Regards, Rentivo")

					if err != nil {
						return nil, err
					}

					var output bytes.Buffer

					err = t.Execute(&output, map[string]string{
						"BuildId":  buildId,
						"SiteName": pipeline.Site.Description,
					})

					if err != nil {
						return nil, err
					}

					err = notify.SendToAllUsers(
						fmt.Sprintf("%s - Build Complete", pipeline.Site.Description),
						output.String(),
						pipeline.Site.AccountId,
					)

					if err != nil {
						return nil, err
					}

					return &types.CommandResult{Command: CmdPreviewBuildNotifyAccountUsers, Output: ""}, nil
				},
			},
		},
	}
}
