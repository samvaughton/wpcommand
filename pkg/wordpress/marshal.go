package wordpress

import (
	"encoding/json"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
)

func UserListFromJson(data string) ([]types.WpUser, error) {
	var list []types.WpUser

	err := json.Unmarshal([]byte(data), &list)

	if err != nil {
		log.Errorf("failed decoding user json list: %s", err)

		return []types.WpUser{}, err
	}

	return list, nil
}

func PluginListFromJson(data string) ([]types.WpPlugin, error) {
	var list []types.WpPlugin

	err := json.Unmarshal([]byte(data), &list)

	if err != nil {
		log.Errorf("failed decoding plugin json list: %s", err)

		return []types.WpPlugin{}, err
	}

	return list, nil
}

func ThemeListFromJson(data string) ([]types.WpTheme, error) {
	var list []types.WpTheme

	err := json.Unmarshal([]byte(data), &list)

	if err != nil {
		log.Errorf("failed decoding theme json list: %s", err)

		return []types.WpTheme{}, err
	}

	return list, nil
}

func PostListFromJson(data string) ([]types.WpPost, error) {
	var list []types.WpPost

	err := json.Unmarshal([]byte(data), &list)

	if err != nil {
		log.Errorf("failed decoding post json list: %s", err)

		return []types.WpPost{}, err
	}

	return list, nil
}
