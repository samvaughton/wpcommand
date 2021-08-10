package pipeline

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/samvaughton/wpcommand/v2/pkg/config"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	"github.com/samvaughton/wpcommand/v2/pkg/wordpress"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

const CmdWpThemesSync = "themes-sync"
const CmdWpThemesStatus = "themes-status"
const CmdWpPluginsSync = "plugins-sync"
const CmdWpPluginsStatus = "plugins-status"
const CmdWpSyncAcfFields = "acf-sync-fields"
const CmdWpSyncLazyblocks = "lazyblocks-sync"
const CmdWpDataUrlTransfers = "data-url-transfers"
const CmdWpUpdateSiteConfig = "update-site-config"

var CommandRegistry = map[string]func(site *types.Site) SiteCommand{
	CmdWpThemesSync:       GetThemesSyncCommand,
	CmdWpThemesStatus:     GetThemesStatusCommand,
	CmdWpPluginsSync:      GetPluginsSyncCommand,
	CmdWpPluginsStatus:    GetPluginsStatusCommand,
	CmdWpSyncAcfFields:    GetAcfSyncFieldsCommand,
	CmdWpSyncLazyblocks:   GetLazyBlocksSyncCommand,
	CmdWpDataUrlTransfers: GetWpDataUrlTransferCommand,
	CmdWpUpdateSiteConfig: GetWpUpdateSiteConfigCommand,
}

func CommandExists(key string) bool {
	_, exists := CommandRegistry[key]

	return exists
}

func GetThemesSyncCommand(site *types.Site) SiteCommand {
	return &WrappedCommand{
		Name: CmdWpThemesSync,
		Wrapped: func(pipeline *SiteCommandPipeline) (*types.CommandResult, error) {
			actionSet, err := wordpress.GetSiteThemeStatuses(site.Id, pipeline.Executor)

			if err != nil {
				return nil, errors.New("failed to pull theme list, cannot sync")
			}

			// Install and upgrade themes first, as the theme we are installing needs to be activated
			wordpress.RunThemeActionOnSet(pipeline.Executor, actionSet, types.ThemeActionEnum.Install)
			wordpress.RunThemeActionOnSet(pipeline.Executor, actionSet, types.ThemeActionEnum.Upgrade)

			// we need to reverse the order of the uninstall/downgrade
			// in case themes error out if the core one is removed first
			// causing PHP errors in some themes that dont handle errors very well
			wordpress.RunThemeActionOnSet(pipeline.Executor, wordpress.ReverseThemeActionItemOrder(actionSet), types.ThemeActionEnum.Uninstall)
			wordpress.RunThemeActionOnSet(pipeline.Executor, wordpress.ReverseThemeActionItemOrder(actionSet), types.ThemeActionEnum.Downgrade)

			return nil, nil
		},
	}
}

func GetThemesStatusCommand(site *types.Site) SiteCommand {
	return &WrappedCommand{
		Name: CmdWpThemesStatus,
		Wrapped: func(pipeline *SiteCommandPipeline) (*types.CommandResult, error) {
			actionSet, err := wordpress.GetSiteThemeStatuses(site.Id, pipeline.Executor)

			if err != nil {
				return nil, errors.New("failed to pull theme list, cannot retrieve status")
			}

			for _, action := range actionSet {
				log.WithFields(log.Fields{
					"site":    site.Key,
					"theme":   action.Object.Name,
					"version": action.Object.Version,
					"action":  action.Action,
				}).Info("theme status")
			}

			return nil, nil
		},
	}
}

func GetPluginsSyncCommand(site *types.Site) SiteCommand {
	return &WrappedCommand{
		Name: CmdWpPluginsSync,
		Wrapped: func(pipeline *SiteCommandPipeline) (*types.CommandResult, error) {
			actionSet, err := wordpress.GetSitePluginStatuses(site.Id, pipeline.Executor)

			if err != nil {
				return nil, errors.New("failed to pull plugin list, cannot sync")
			}

			wordpress.RunPluginActionOnSet(pipeline.Executor, wordpress.ReversePluginActionItemOrder(actionSet), types.PluginActionEnum.Uninstall)
			wordpress.RunPluginActionOnSet(pipeline.Executor, wordpress.ReversePluginActionItemOrder(actionSet), types.PluginActionEnum.Downgrade)

			wordpress.RunPluginActionOnSet(pipeline.Executor, actionSet, types.PluginActionEnum.Install)
			wordpress.RunPluginActionOnSet(pipeline.Executor, actionSet, types.PluginActionEnum.Upgrade)

			return nil, nil
		},
	}
}

func GetPluginsStatusCommand(site *types.Site) SiteCommand {
	return &WrappedCommand{
		Name: CmdWpPluginsStatus,
		Wrapped: func(pipeline *SiteCommandPipeline) (*types.CommandResult, error) {
			actionSet, err := wordpress.GetSitePluginStatuses(site.Id, pipeline.Executor)

			if err != nil {
				return nil, errors.New("failed to pull plugin list, cannot retrieve status")
			}

			for _, action := range actionSet {
				log.WithFields(log.Fields{
					"site":    site.Key,
					"theme":   action.Object.Name,
					"version": action.Object.Version,
					"action":  action.Action,
				}).Info("plugin status")
			}

			return nil, nil
		},
	}
}

func GetWpDataUrlTransferCommand(site *types.Site) SiteCommand {
	return &WrappedCommand{
		Name: CmdWpDataUrlTransfers,
		Wrapped: func(pipeline *SiteCommandPipeline) (*types.CommandResult, error) {
			for _, item := range pipeline.Config.Wordpress.DataUrls {
				fullUrl := fmt.Sprintf("%s%s", config.Config.Wordpress.StoreRoot, item.Path)

				downloadCmd := fmt.Sprintf("curl %s --output /opt/bitnami/%s", fullUrl, item.Name)
				_, err := pipeline.Executor.ExecuteCommand([]string{downloadCmd})

				if err != nil {
					return nil, errors.New(fmt.Sprintf("could not transfer file to the pod - all files need to be transferred before starting setup. %s", err))
				}
			}

			return nil, nil
		},
	}
}

func GetAcfSyncFieldsCommand(site *types.Site) SiteCommand {
	return &SimplePipelineCommand{
		Name: CmdWpSyncAcfFields,
		Commands: []SiteCommand{
			// delete ACF field groups, ready for re-sync
			&DynamicArgsCommand{
				Name: "delete afc field groups",
				GetArgs: func(pipeline *SiteCommandPipeline) ([]string, error) {
					// now we need to delete all lazy blocks and import the new ones
					// Delete any default posts
					var ids []string

					acfFieldGroups, err := wordpress.GetSiteAcfFieldGroups(pipeline.Executor)

					if err != nil {
						return []string{}, errors.New("failed to get acf field group list, cannot sync")
					}

					for _, block := range acfFieldGroups {
						ids = append(ids, strconv.Itoa(block.Id))
					}

					if len(ids) > 0 {
						return []string{fmt.Sprintf("wp post delete %s --force", strings.Join(ids, " "))}, nil
					}

					return []string{}, nil // no execution
				},
			},
			&SimpleCommand{
				Name: "sync afc field groups",
				Args: []string{"wp eval-file /opt/bitnami/eval-index.php acf-field-group-sync"},
			},
		},
	}
}

func GetLazyBlocksSyncCommand(site *types.Site) SiteCommand {
	return &SimplePipelineCommand{
		Name: CmdWpSyncLazyblocks,
		Commands: []SiteCommand{
			// delete old lazy blocks
			&DynamicArgsCommand{
				Name: "delete existing lazyblocks",
				GetArgs: func(pipeline *SiteCommandPipeline) ([]string, error) {
					// now we need to delete all lazy blocks and import the new ones
					// Delete any default posts
					var lazyIds []string
					lazyblocks, err := wordpress.GetSiteLazyblocksPosts(pipeline.Executor)

					if err != nil {
						return []string{}, errors.New("failed to get lazy block list, cannot sync")
					}

					for _, block := range lazyblocks {
						lazyIds = append(lazyIds, strconv.Itoa(block.Id))
					}

					if len(lazyIds) > 0 {
						return []string{fmt.Sprintf("wp post delete %s --force", strings.Join(lazyIds, " "))}, nil
					}

					return []string{}, nil
				},
			},
			// import lazy blocks
			&SimpleCommand{
				Name: "import new lazyblocks",
				Args: []string{"wp eval-file /opt/bitnami/eval-index.php lazyblocks-import"},
			},
		},
	}
}

func GetWpUpdateSiteConfigCommand(site *types.Site) SiteCommand {
	return &SimplePipelineCommand{
		Name:     CmdWpUpdateSiteConfig,
		Commands: []SiteCommand{
			/*&WrappedCommand{
				Name: "rentivo siteConfig.js transfer",
				Wrapped: func(pipeline *SiteCommandPipeline) *types.CommandResult {
					downloadCmd := fmt.Sprintf("curl %s --output /opt/bitnami/siteConfig.json", "localhost")
					result := pipeline.Executor.ExecuteCommand([]string{downloadCmd})

					if result.HasError() {
						log.Fatal("could not transfer file to the pod - all files need to be transferred before starting setup")
					}

					return nil
				},
			},
			&SimpleCommand{
				Name: "set rentivo siteConfig.js",
				Args: []string{"wp eval-file /opt/bitnami/eval-index.php rentivo-simba-update-site-config"},
			},*/
		},
	}
}
