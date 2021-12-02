package registry

import (
	"fmt"
	"github.com/samvaughton/wpcommand/v2/pkg/pipeline"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
)

func GetSetDefaultOptionsCommand(site *types.Site) pipeline.SiteCommand {
	return &pipeline.SimplePipelineCommand{
		Name: CmdWpSetDefaultOptions,
		Commands: []pipeline.SiteCommand{
			&pipeline.SimpleCommand{
				Name: fmt.Sprintf("%s.%s", CmdWpSetDefaultOptions, "show_on_front=page"),
				Args: []string{"wp option update show_on_front page"},
			},
			&pipeline.SimpleCommand{
				Name: fmt.Sprintf("%s.%s", CmdWpSetDefaultOptions, "permalink_structure=%postname%"),
				Args: []string{"wp option update permalink_structure \"/%postname%/\""},
			},
			&pipeline.SimpleCommand{
				Name: fmt.Sprintf("%s.%s", CmdWpSetDefaultOptions, "wpcmd_site_access_token=***"),
				Args: []string{fmt.Sprintf("wp option update wpcmd_site_access_token \"%s\"", site.AccessToken)},
			},
		},
	}
}
