package registry

import (
	"fmt"
	"github.com/samvaughton/wpcommand/v2/pkg/db"
	"github.com/samvaughton/wpcommand/v2/pkg/pipeline"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	"github.com/samvaughton/wpcommand/v2/pkg/wordpress"
)

func GetUserSetupCommand(site *types.Site, config map[string]interface{}) pipeline.SiteCommand {
	return &pipeline.SimplePipelineCommand{
		Name:              CmdWpSiteUserSetup,
		ErrorIsSuccessful: true, // if no user is found, it errors and then we create
		RunPreCheckCommand: &pipeline.SimpleCommand{
			Name: CmdWpSiteUserCheck,
			Args: []string{fmt.Sprintf("wp user get %s --format=json", site.SiteUsername)},
		},
		Commands: []pipeline.SiteCommand{
			&pipeline.SimpleCommand{
				Name: CmdWpSiteUserCreate,
				Args: []string{fmt.Sprintf("wp user create %s %s --user_pass=%s --role=administrator", site.SiteUsername, site.SiteEmail, site.SitePassword)},
			},
		},
	}
}

func GetUserSyncCommand(site *types.Site, config map[string]interface{}) pipeline.SiteCommand {
	return &pipeline.WrappedCommand{
		Name: CmdWpUserList,
		Wrapped: func(pipeline *pipeline.SiteCommandPipeline) (*types.CommandResult, error) {
			list, err := wordpress.GetSiteUserList(site.Id, pipeline.Executor)

			if err != nil {
				return nil, err
			}

			cached, err := site.GetWpCachedData()

			if err != nil {
				return nil, err
			}

			cached.UserList = list

			err = site.SetWpCachedData(&cached)

			if err != nil {
				return nil, err
			}

			db.SiteUpdate(site)

			return &types.CommandResult{Data: list}, nil
		},
	}
}
