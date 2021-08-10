package wordpress

import (
	"encoding/json"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
)

func PluginListFromJson(data string) ([]types.Plugin, error) {
	var list []types.Plugin

	err := json.Unmarshal([]byte(data), &list)

	if err != nil {
		log.Errorf("failed decoding plugin json list: %s", err)

		return []types.Plugin{}, err
	}

	return list, nil
}

func ThemeListFromJson(data string) ([]types.Theme, error) {
	var list []types.Theme

	err := json.Unmarshal([]byte(data), &list)

	if err != nil {
		log.Errorf("failed decoding theme json list: %s", err)

		return []types.Theme{}, err
	}

	return list, nil
}

func PostListFromJson(data string) ([]types.Post, error) {
	var list []types.Post

	err := json.Unmarshal([]byte(data), &list)

	if err != nil {
		log.Errorf("failed decoding post json list: %s", err)

		return []types.Post{}, err
	}

	return list, nil
}
