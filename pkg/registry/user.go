package registry

import (
	"fmt"
	"github.com/samvaughton/wpcommand/v2/pkg/pipeline"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
)

func GetUserSetupCommand(site *types.Site) pipeline.SiteCommand {
	return &pipeline.SimplePipelineCommand{
		Name:              CmdWpUserSetup,
		ErrorIsSuccessful: true, // if no user is found, it errors and then we create
		RunPreCheckCommand: &pipeline.SimpleCommand{
			Name: CmdWpUserCheck,
			Args: []string{fmt.Sprintf("wp user get %s --format=json", site.SiteUsername)},
		},
		Commands: []pipeline.SiteCommand{
			&pipeline.SimpleCommand{
				Name: CmdWpUserCreate,
				Args: []string{fmt.Sprintf("wp user create %s %s --user_pass=%s --role=administrator", site.SiteUsername, site.SiteEmail, site.SitePassword)},
			},
		},
	}
}
