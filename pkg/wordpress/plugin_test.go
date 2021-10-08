package wordpress

import (
	"github.com/google/uuid"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

const pluginStatusOutput = `[
	{"name":"acf-to-rest-api","status":"active","update":"none","version":"3.3.2"},
	{"name":"acf-code-field","status":"active","update":"none","version":"1.8"},
	{"name":"acf-map-bbox","status":"active","update":"none","version":"1.0.0"},
	{"name":"duplicate-page","status":"active","update":"available","version":"4.4"}
]`

func TestGetSitePluginStatuses(t *testing.T) {
	latestObjs := []types.ObjectBlueprint{
		{ // verify action NONE
			Id:         1,
			RevisionId: 1,
			ExactName:  "acf-to-rest-api",
			Uuid:       uuid.New().String(),
			Type:       types.ObjectBlueprintTypePlugin,
			SetOrder:   10,
			Active:     true,
			Version:    "3.3.2",
		},
		{ // verify action DOWNGRADE
			Id:         2,
			RevisionId: 1,
			ExactName:  "acf-code-field",
			Uuid:       uuid.New().String(),
			Type:       types.ObjectBlueprintTypePlugin,
			SetOrder:   20,
			Active:     true,
			Version:    "1.2",
		},
		{ // verify action UPGRADE
			Id:         3,
			RevisionId: 1,
			ExactName:  "acf-map-bbox",
			Uuid:       uuid.New().String(),
			Type:       types.ObjectBlueprintTypePlugin,
			SetOrder:   30,
			Active:     true,
			Version:    "1.1.34",
		},
		{ // verify action INSTALL
			Id:         4,
			RevisionId: 1,
			ExactName:  "advanced-custom-fields-pro",
			Uuid:       uuid.New().String(),
			Type:       types.ObjectBlueprintTypePlugin,
			SetOrder:   40,
			Active:     true,
			Version:    "5.9.5",
		},
		// omit duplicate page to mimic UNINSTALL
	}

	getActionByExactName := func(exactName string, set types.PluginActionSet) types.PluginAction {
		for _, item := range set.Items {
			if item.Name == exactName {
				return item.Action
			}
		}

		return ""
	}

	t.Run("GetSitePluginStatuses", func(t *testing.T) {
		result, err := GetSitePluginStatuses(newTestExecutor(pluginStatusOutput), latestObjs)

		assert.Nil(t, err)

		assert.Equal(t, types.PluginActionEnum.None, getActionByExactName("acf-to-rest-api", result))
		assert.Equal(t, types.PluginActionEnum.Downgrade, getActionByExactName("acf-code-field", result))
		assert.Equal(t, types.PluginActionEnum.Upgrade, getActionByExactName("acf-map-bbox", result))
		assert.Equal(t, types.PluginActionEnum.Install, getActionByExactName("advanced-custom-fields-pro", result))
		assert.Equal(t, types.PluginActionEnum.Uninstall, getActionByExactName("duplicate-page", result))
	})
}
