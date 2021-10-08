package registry

import (
	"errors"
	"fmt"
	"github.com/samvaughton/wpcommand/v2/pkg/db"
	"github.com/samvaughton/wpcommand/v2/pkg/pipeline"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	"github.com/samvaughton/wpcommand/v2/pkg/wordpress"
)

func GetThemesStatusCommand(site *types.Site) pipeline.SiteCommand {
	return &pipeline.WrappedCommand{
		Name: CmdWpThemesStatus,
		Wrapped: func(pipeline *pipeline.SiteCommandPipeline) (*types.CommandResult, error) {
			latestObjs := db.GetLatestObjectBlueprintsForSiteAndType(site.Id, types.ObjectBlueprintTypeTheme)
			actionSet, err := wordpress.GetSiteThemeStatuses(pipeline.Executor, latestObjs)

			if err != nil {
				return nil, err
			}

			return &types.CommandResult{Data: actionSet}, nil
		},
	}
}

func GetThemesSyncCommand(site *types.Site) pipeline.SiteCommand {
	return &pipeline.SimplePipelineCommand{
		Name: CmdWpThemesSync,
		Commands: []pipeline.SiteCommand{
			GetThemesStatusCommand(site),
			&pipeline.WrappedCommand{
				Name: fmt.Sprintf("%s.%s", CmdWpThemesSync, "execute"),
				Wrapped: func(pipeline *pipeline.SiteCommandPipeline) (*types.CommandResult, error) {
					actionSet, valid := pipeline.PreviousResult.Data.(types.ThemeActionSet)

					if valid == false {
						return nil, errors.New("could not load in action set from previous command")
					}

					// install and activate first for theme so we can deactivate the old one
					err := wordpress.RunThemeActionOnSet(pipeline.Executor, actionSet.Items, types.ThemeActionEnum.Install)

					if err != nil {
						return nil, err
					}

					err = wordpress.RunThemeActionOnSet(pipeline.Executor, actionSet.Items, types.ThemeActionEnum.Upgrade)

					if err != nil {
						return nil, err
					}

					err = wordpress.RunThemeActionOnSet(pipeline.Executor, wordpress.ReverseThemeActionItemOrder(actionSet.Items), types.ThemeActionEnum.Uninstall)

					if err != nil {
						return nil, err
					}

					err = wordpress.RunThemeActionOnSet(pipeline.Executor, wordpress.ReverseThemeActionItemOrder(actionSet.Items), types.ThemeActionEnum.Downgrade)

					if err != nil {
						return nil, err
					}

					return nil, nil
				},
			},
		},
	}
}
