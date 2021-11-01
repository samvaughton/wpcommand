package registry

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/samvaughton/wpcommand/v2/pkg/db"
	"github.com/samvaughton/wpcommand/v2/pkg/pipeline"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	"github.com/samvaughton/wpcommand/v2/pkg/wordpress"
	log "github.com/sirupsen/logrus"
)

func GetPluginsStatusCommand(site *types.Site) pipeline.SiteCommand {
	return &pipeline.WrappedCommand{
		Name: CmdWpPluginsStatus,
		Wrapped: func(pipeline *pipeline.SiteCommandPipeline) (*types.CommandResult, error) {
			latestObjs := db.GetLatestObjectBlueprintsForSiteAndType(site.Id, types.ObjectBlueprintTypePlugin)
			actionSet, err := wordpress.GetSitePluginStatuses(pipeline.Executor, latestObjs)

			if err != nil {
				return nil, err
			}

			for _, action := range actionSet.Items {
				ver := ""
				if action.Object != nil {
					ver = action.Object.Version
				}

				log.WithFields(log.Fields{
					"site":    site.Key,
					"plugin":  action.Name,
					"version": ver,
					"action":  action.Action,
				}).Info("plugin status")
			}

			return &types.CommandResult{Data: actionSet}, nil
		},
	}
}

func GetPluginsSyncCommand(site *types.Site) pipeline.SiteCommand {
	return &pipeline.SimplePipelineCommand{
		Name: CmdWpPluginsSync,
		Commands: []pipeline.SiteCommand{
			GetPluginsStatusCommand(site),
			&pipeline.WrappedCommand{
				Name: fmt.Sprintf("%s.%s", CmdWpPluginsSync, "execute"),
				Wrapped: func(pipeline *pipeline.SiteCommandPipeline) (*types.CommandResult, error) {
					actionSet, valid := pipeline.PreviousResult.Data.(types.PluginActionSet)

					if valid == false {
						return nil, errors.New("could not load in action set from previous command")
					}

					err := wordpress.RunPluginActionOnSet(pipeline.Executor, wordpress.ReversePluginActionItemOrder(actionSet.Items), types.PluginActionEnum.Uninstall)

					if err != nil {
						return nil, err
					}

					err = wordpress.RunPluginActionOnSet(pipeline.Executor, wordpress.ReversePluginActionItemOrder(actionSet.Items), types.PluginActionEnum.Downgrade)

					if err != nil {
						return nil, err
					}

					err = wordpress.RunPluginActionOnSet(pipeline.Executor, actionSet.Items, types.PluginActionEnum.Install)

					if err != nil {
						return nil, err
					}

					err = wordpress.RunPluginActionOnSet(pipeline.Executor, actionSet.Items, types.PluginActionEnum.Upgrade)

					if err != nil {
						return nil, err
					}

					return nil, nil
				},
			},
		},
	}
}
