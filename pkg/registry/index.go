package registry

import (
	"github.com/samvaughton/wpcommand/v2/pkg/db"
	"github.com/samvaughton/wpcommand/v2/pkg/pipeline"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
	"strings"
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

func CreateDefaultCommands() {
	commands, err := db.CommandsGetDefault()

	if err != nil {
		log.Error(err)

		return
	}

	if len(commands) > 0 {
		log.Info("default commands exist, skipping")

		return
	}

	for key, _ := range CommandRegistry {
		description := strings.Title(strings.ReplaceAll(key, "-", " "))

		cmd, err := db.CommandCreateDefault(description, key, types.CommandTypeWpBuiltIn)

		if err != nil {
			log.Error(err)

			continue
		}

		log.Infof("command %s creeated", cmd.Key)
	}
}
