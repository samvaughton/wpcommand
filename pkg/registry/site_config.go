package registry

import (
	"fmt"
	"github.com/samvaughton/wpcommand/v2/pkg/pipeline"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
)

func GetWpUpdateSiteConfigCommand(site *types.Site) pipeline.SiteCommand {
	return &pipeline.SimplePipelineCommand{
		Name: CmdWpUpdateSiteConfig,
		Commands: []pipeline.SiteCommand{
			&pipeline.WrappedCommand{
				Name: fmt.Sprintf("%s.%s", CmdWpUpdateSiteConfig, "download"),
				Wrapped: func(pipeline *pipeline.SiteCommandPipeline) (*types.CommandResult, error) {
					// this actually needs to be a live endpoint
					scUrl := fmt.Sprintf("http://00e9a9271f0d.ngrok.io/public/site/%s/config", site.Uuid)

					downloadCmd := fmt.Sprintf("curl %s --output /opt/bitnami/siteConfig.json", scUrl)
					result, err := pipeline.Executor.ExecuteCommand([]string{downloadCmd})

					if err != nil {
						return nil, err
					}

					return result, nil
				},
			},
			&pipeline.SimpleCommand{
				Name: fmt.Sprintf("%s.%s", CmdWpUpdateSiteConfig, "apply"),
				Args: []string{"wp eval-file /opt/bitnami/eval-index.php rentivo-simba-update-site-config"},
			},
		},
	}
}
