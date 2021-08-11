package registry

import (
	"github.com/samvaughton/wpcommand/v2/pkg/pipeline"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
)

const CmdWpSetupFreshInstall = "setup-fresh-install"

const CmdWpThemesSync = "themes-sync"
const CmdWpThemesStatus = "themes-status"

const CmdWpPluginsSync = "plugins-sync"
const CmdWpPluginsStatus = "plugins-status"

const CmdWpPolylangSetup = "polylang-setup"
const CmdWpPolylangCliInstall = "polylang-cli-install"
const CmdWpPolylangConfigure = "polylang-configure"

const CmdWpUserSetup = "user-setup"
const CmdWpUserCheck = "user-check"
const CmdWpUserCreate = "user-create"

const CmdWpHousecleaning = "wp-housecleaning"

const CmdWpSetDefaultOptions = "wp-set-default-options"

const CmdWpSyncAcfFields = "acf-sync-fields"
const CmdWpSyncLazyblocks = "lazyblocks-sync"
const CmdWpDataUrlTransfers = "data-url-transfers"
const CmdWpUpdateSiteConfig = "update-site-config"

var CommandRegistry = map[string]func(site *types.Site) pipeline.SiteCommand{
	CmdWpThemesSync:       GetThemesSyncCommand,
	CmdWpThemesStatus:     GetThemesStatusCommand,
	CmdWpPluginsSync:      GetPluginsSyncCommand,
	CmdWpPluginsStatus:    GetPluginsStatusCommand,
	CmdWpSyncAcfFields:    GetAcfSyncFieldsCommand,
	CmdWpSyncLazyblocks:   GetLazyblocksSyncCommand,
	CmdWpDataUrlTransfers: GetWpDataUrlTransferCommand,

	CmdWpUserSetup:         GetUserSetupCommand,
	CmdWpHousecleaning:     GetHousecleaningCommand,
	CmdWpSetDefaultOptions: GetSetDefaultOptionsCommand,
	CmdWpPolylangSetup:     GetPolylangSetupCommand,

	CmdWpSetupFreshInstall: GetWpSetupFreshInstallCommand,
}

func CommandExists(key string) bool {
	_, exists := CommandRegistry[key]

	return exists
}
