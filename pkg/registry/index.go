package registry

import (
	"github.com/samvaughton/wpcommand/v2/pkg/db"
	"github.com/samvaughton/wpcommand/v2/pkg/pipeline"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
	"strings"
)

const CmdWpSetupFreshInstall = "setup-fresh-install"

const CmdWpCliLoginInstall = "wp-cli-login-install"
const CmdWpCliLoginCreateLoginUrl = "wp-cli-login-create-login-url"

const CmdWpThemesSync = "themes-sync"
const CmdWpThemesStatus = "themes-status"

const CmdWpPluginsSync = "plugins-sync"
const CmdWpPluginsStatus = "plugins-status"

const CmdWpPolylangSetup = "polylang-setup"
const CmdWpPolylangCliInstall = "polylang-cli-install"
const CmdWpPolylangConfigure = "polylang-configure"

const CmdWpSiteUserSetup = "site-user-setup"
const CmdWpSiteUserCheck = "site-user-check"
const CmdWpSiteUserCreate = "site-user-create"

const CmdWpSiteUserSync = "wp-user-sync"
const CmdWpUserList = "wp-user-list"

// currently done manually
//const CmdWpUserDelete = "wp-user-delete"
//const CmdWpUserCreate = "wp-user-create"
//const CmdWpUserUpdate = "wp-user-update"

const CmdWpHousecleaning = "wp-housecleaning"

const CmdWpSetDefaultOptions = "wp-set-default-options"

const CmdWpSyncAcfFields = "acf-sync-fields"
const CmdWpSyncLazyblocks = "lazyblocks-sync"
const CmdWpDataUrlTransfers = "data-url-transfers"
const CmdWpUpdateSiteConfig = "update-site-config"

const CmdHttpCall = "http-call"

const CmdPreviewBuild = "preview-build"
const CmdPreviewBuildRunGithubWorkflow = "preview-build-github-workflow"
const CmdPreviewBuildRunDockerRegistryWorkflow = "preview-build-docker-registry"
const CmdPreviewBuildKubernetesDeploy = "preview-build-kubernetes-deploy"
const CmdPreviewBuildNotifyAccountUsers = "preview-build-notify-users"

const CmdBuildRelease = "build-release"
const CmdBuildReleaseRunGithubWorkflow = "build-release-run-github-workflow"
const CmdBuildReleaseCheckDeploymentStatus = "build-release-check-deployment-status"
const CmdBuildReleaseNotifyAccountUsers = "build-release-notify-users"

var BuiltInCommandRegistry = map[string]func(site *types.Site) pipeline.SiteCommand{
	CmdWpThemesSync:       GetThemesSyncCommand,
	CmdWpThemesStatus:     GetThemesStatusCommand,
	CmdWpPluginsSync:      GetPluginsSyncCommand,
	CmdWpPluginsStatus:    GetPluginsStatusCommand,
	CmdWpSyncAcfFields:    GetAcfSyncFieldsCommand,
	CmdWpSyncLazyblocks:   GetLazyblocksSyncCommand,
	CmdWpDataUrlTransfers: GetWpDataUrlTransferCommand,

	CmdWpSiteUserSetup: GetUserSetupCommand,
	CmdWpSiteUserSync:  GetUserSyncCommand,

	CmdWpHousecleaning:     GetHousecleaningCommand,
	CmdWpSetDefaultOptions: GetSetDefaultOptionsCommand,
	CmdWpPolylangSetup:     GetPolylangSetupCommand,

	CmdWpSetupFreshInstall: GetWpSetupFreshInstallCommand,
	CmdWpCliLoginInstall:   GetWpCliLoginInstallCommand,
}

func CommandExists(key string) bool {
	_, exists := BuiltInCommandRegistry[key]

	return exists
}

func CreateDefaultCommands() {
	for key, _ := range BuiltInCommandRegistry {
		existingCmd, err := db.CommandGetByKey(key)

		if existingCmd != nil {
			continue
		}

		description := strings.Title(strings.ReplaceAll(key, "-", " "))

		cmd, err := db.CommandCreateDefault(description, key, types.CommandTypeWpBuiltIn)

		if err != nil {
			log.Error(err)

			continue
		}

		log.Infof("command %s created", cmd.Key)
	}
}
