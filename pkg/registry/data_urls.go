package registry

import (
	"fmt"
	appConfig "github.com/samvaughton/wpcommand/v2/pkg/config"
	"github.com/samvaughton/wpcommand/v2/pkg/pipeline"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
)

func GetWpDataUrlTransferCommand(site *types.Site) pipeline.SiteCommand {
	return &pipeline.WrappedCommand{
		Name: CmdWpDataUrlTransfers,
		Wrapped: func(pipeline *pipeline.SiteCommandPipeline) (*types.CommandResult, error) {
			for _, item := range pipeline.Config.Wordpress.DataUrls {
				fullUrl := fmt.Sprintf("%s%s", appConfig.Config.Wordpress.StoreRoot, item.Path)

				downloadCmd := fmt.Sprintf("curl %s --output /opt/bitnami/%s", fullUrl, item.Name)
				_, err := pipeline.Executor.ExecuteCommand([]string{downloadCmd})

				if err != nil {
					return nil, err
				}
			}

			return nil, nil
		},
	}
}
