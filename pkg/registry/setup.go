package registry

import (
	"github.com/samvaughton/wpcommand/v2/pkg/pipeline"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
)

func GetWpSetupFreshInstallCommand(site *types.Site, config map[string]interface{}) pipeline.SiteCommand {
	var commands []pipeline.SiteCommand

	commands = append(commands, GetPluginsSyncCommand(site, config))
	commands = append(commands, GetThemesSyncCommand(site, config))
	commands = append(commands, GetWpDataUrlTransferCommand(site, config))
	commands = append(commands, GetPolylangSetupCommand(site, config))

	if site.SiteUsername != "" && site.SitePassword != "" && site.SiteEmail != "" {
		commands = append(commands, GetUserSetupCommand(site, config))
	}

	commands = append(commands, GetHousecleaningCommand(site, config))
	commands = append(commands, GetLazyblocksSyncForFreshInstallCommand(site, config)) // doesnt include data url download
	commands = append(commands, GetAcfSyncFieldsCommand(site, config))
	commands = append(commands, GetSetDefaultOptionsCommand(site, config))
	commands = append(commands, GetWpUpdateSiteConfigCommand(site))

	return &pipeline.SimplePipelineCommand{
		Name:     CmdWpSetupFreshInstall,
		Commands: commands,
	}
}
