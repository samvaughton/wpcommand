package registry

import (
	"github.com/samvaughton/wpcommand/v2/pkg/pipeline"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
)

func GetWpSetupFreshInstallCommand(site *types.Site) pipeline.SiteCommand {
	var commands []pipeline.SiteCommand

	commands = append(commands, GetPluginsSyncCommand(site))
	commands = append(commands, GetThemesSyncCommand(site))
	commands = append(commands, GetWpDataUrlTransferCommand(site))
	commands = append(commands, GetPolylangSetupCommand(site))

	if site.SiteUsername != "" && site.SitePassword != "" && site.SiteEmail != "" {
		commands = append(commands, GetUserSetupCommand(site))
	}

	commands = append(commands, GetHousecleaningCommand(site))
	commands = append(commands, GetLazyblocksSyncCommand(site))
	commands = append(commands, GetAcfSyncFieldsCommand(site))
	commands = append(commands, GetSetDefaultOptionsCommand(site))
	commands = append(commands, GetWpUpdateSiteConfigCommand(site))

	return &pipeline.SimplePipelineCommand{
		Name:     CmdWpSetupFreshInstall,
		Commands: commands,
	}
}
