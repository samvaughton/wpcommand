package registry

import (
	"fmt"
	"github.com/samvaughton/wpcommand/v2/pkg/pipeline"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
)

func GetPolylangSetupCommand(site *types.Site) pipeline.SiteCommand {
	return &pipeline.SimplePipelineCommand{
		Name: CmdWpPolylangSetup,
		Commands: []pipeline.SiteCommand{
			GetPolylangCliInstallCommand(site),
			GetPolylangConfigureCommand(site),
		},
	}
}

func GetPolylangCliInstallCommand(site *types.Site) pipeline.SiteCommand {
	return &pipeline.SimpleCommand{
		Name: CmdWpPolylangCliInstall,
		Args: []string{"wp package install https://github.com/aptenex/polylang-cli/archive/refs/heads/master.zip --insecure"},
	}
}

func GetPolylangConfigureCommand(site *types.Site) pipeline.SiteCommand {
	return &pipeline.SimplePipelineCommand{
		Name:              CmdWpPolylangConfigure,
		ErrorIsSuccessful: false,
		RunPreCheckCommand: &pipeline.SimpleCommand{
			Args: []string{"wp package path aptenex/polylang-cli"},
		},
		Commands: []pipeline.SiteCommand{
			// create lang
			&pipeline.SimpleCommand{
				Name: fmt.Sprintf("%s.%s", CmdWpPolylangConfigure, "create-lang"),
				Args: []string{fmt.Sprintf("wp pll lang create %s %s %s --flag=us", "English", "en", "en_US")},
			},
			// defaults
			&pipeline.SimpleCommand{
				Name: fmt.Sprintf("%s.%s", CmdWpPolylangConfigure, "set-default-lang=en"),
				Args: []string{fmt.Sprintf("wp pll option default %s", "en")},
			},
			// enable post types
			&pipeline.SimpleCommand{
				Name: fmt.Sprintf("%s.%s", CmdWpPolylangConfigure, "post-types-enable=collections,experiences"),
				Args: []string{fmt.Sprintf("wp pll post-type enable %s", "collections,experiences")},
			},
			// enable taxonomies
			&pipeline.SimpleCommand{
				Name: fmt.Sprintf("%s.%s", CmdWpPolylangConfigure, "taxonomy-enable=experience_taxonomy,experience_tags"),
				Args: []string{fmt.Sprintf("wp pll taxonomy enable %s", "experience_taxonomy,experience_tags")},
			},
		},
	}
}
