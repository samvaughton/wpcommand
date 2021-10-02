package wordpress

import (
	"github.com/samvaughton/wpcommand/v2/pkg/execution"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
)

func GetSiteLazyblocksPosts(executor execution.CommandExecutor) ([]types.WpPost, error) {
	result, err := executor.ExecuteCommand([]string{"wp post list --post_type=lazyblocks --format=json"})

	if err != nil {
		return []types.WpPost{}, err
	}

	return PostListFromJson(result.Output)
}

func GetSitePostAndPages(executor execution.CommandExecutor) ([]types.WpPost, error) {
	result, err := executor.ExecuteCommand([]string{"wp post list --post_type=page,post --format=json"})

	if err != nil {
		return []types.WpPost{}, err
	}

	return PostListFromJson(result.Output)
}

func GetSiteAcfFieldGroups(executor execution.CommandExecutor) ([]types.WpPost, error) {
	result, err := executor.ExecuteCommand([]string{"wp post list --post_type=acf-field-group --format=json"})

	if err != nil {
		return []types.WpPost{}, err
	}

	return PostListFromJson(result.Output)
}
