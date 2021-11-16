package registry

import (
	"github.com/samvaughton/wpcommand/v2/pkg/pipeline"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
)

func GetWpCliLoginInstallCommand(site *types.Site) pipeline.SiteCommand {
	return &pipeline.SimplePipelineCommand{Commands: []pipeline.SiteCommand{
		&pipeline.SimpleCommand{
			Name: CmdWpCliLoginInstall,
			Args: []string{"wp package install aaemnnosttv/wp-cli-login-command"},
		},
		&pipeline.SimpleCommand{
			Name: CmdWpCliLoginInstall,
			Args: []string{"wp login install --activate --yes"},
		},
	}}
}
