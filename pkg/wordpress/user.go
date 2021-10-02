package wordpress

import (
	"github.com/samvaughton/wpcommand/v2/pkg/execution"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
)

func GetSiteUserList(siteId int64, executor execution.CommandExecutor) ([]types.WpUser, error) {
	result, err := executor.ExecuteCommand([]string{"wp user list --format=json --fields=ID,user_login,display_name,user_email,user_registered,roles,user_pass,user_status"})

	if err != nil {
		return []types.WpUser{}, err
	}

	list, err := UserListFromJson(result.Output)

	if err != nil {
		return []types.WpUser{}, err
	}

	return list, nil
}
