package registry

import (
	"fmt"
	"github.com/samvaughton/wpcommand/v2/pkg/pipeline"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
)

func GetWpSetupFreshInstallCommand(site *types.Site) pipeline.SiteCommand {
	var commands []pipeline.SiteCommand

	commands = append(commands, GetPluginsSyncCommand(site))
	commands = append(commands, GetThemesSyncCommand(site))
	commands = append(commands, GetWpCliLoginInstallCommand(site))
	commands = append(commands, GetWpDataUrlTransferCommand(site))
	commands = append(commands, GetPolylangSetupCommand(site))

	if site.SiteUsername != "" && site.SitePassword != "" && site.SiteEmail != "" {
		commands = append(commands, GetUserSetupCommand(site))
	}

	commands = append(commands, GetHousecleaningCommand(site))
	commands = append(commands, GetLazyblocksSyncForFreshInstallCommand(site)) // doesnt include data url download
	commands = append(commands, GetAcfSyncFieldsCommand(site))
	commands = append(commands, GetSetDefaultOptionsCommand(site))
	commands = append(commands, GetWpUpdateSiteConfigCommand(site))

	commands = append(commands, &pipeline.SimpleCommand{
		Name: fmt.Sprintf("%s.%s", CmdWpSetDefaultOptions, "page_on_front=566"),
		Args: []string{"wp option update page_on_front 566"},
	})

	return &pipeline.SimplePipelineCommand{
		Name:     CmdWpSetupFreshInstall,
		Commands: commands,
	}
}
