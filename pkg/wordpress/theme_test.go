package wordpress

import (
	"github.com/google/uuid"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

const themeStatusOutput = `[
	{"name":"twentytwentyone","status":"active","update":"available","version":"1.1"},
	{"name":"examplethemefive","status":"active","update":"available","version":"1.0.4"},
	{"name":"savanna","status":"active","update":"available","version":"2.1.0"}
]`

func TestGetSiteThemeStatuses(t *testing.T) {
	latestObjs := []types.ObjectBlueprint{
		{
			Id:         1,
			RevisionId: 1,
			ExactName:  "twentytwentyone",
			Uuid:       uuid.New().String(),
			Type:       types.ObjectBlueprintTypeTheme,
			SetOrder:   10,
			Active:     true,
			Version:    "1.1",
		},
		{ // verify action DOWNGRADE
			Id:         2,
			RevisionId: 1,
			ExactName:  "savanna",
			Uuid:       uuid.New().String(),
			Type:       types.ObjectBlueprintTypeTheme,
			SetOrder:   20,
			Active:     true,
			Version:    "2.1.2",
		},
	}

	getActionByExactName := func(exactName string, set types.ThemeActionSet) types.ThemeAction {
		for _, item := range set.Items {
			if item.Name == exactName {
				return item.Action
			}
		}

		return ""
	}

	t.Run("GetSiteThemeStatuses", func(t *testing.T) {
		result, err := GetSiteThemeStatuses(newTestExecutor(themeStatusOutput), latestObjs)

		assert.Nil(t, err)

		assert.Equal(t, types.ThemeActionEnum.None, getActionByExactName("twentytwentyone", result))
		assert.Equal(t, types.ThemeActionEnum.Upgrade, getActionByExactName("savanna", result))
		assert.Equal(t, types.ThemeActionEnum.Uninstall, getActionByExactName("examplethemefive", result))
	})
}
