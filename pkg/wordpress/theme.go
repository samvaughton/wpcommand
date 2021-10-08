package wordpress

import (
	"errors"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/samvaughton/wpcommand/v2/pkg/execution"
	"github.com/samvaughton/wpcommand/v2/pkg/object_blueprint"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
	"reflect"
	"sort"
	"time"
)

func GetSiteThemeStatuses(executor execution.CommandExecutor, latestObjs []types.ObjectBlueprint) (types.ThemeActionSet, error) {
	result, err := executor.ExecuteCommand([]string{"wp theme list --format=json"})

	if err != nil {
		return types.ThemeActionSet{}, err
	}

	themeList, err := ThemeListFromJson(result.Output)

	if err != nil {
		return types.ThemeActionSet{}, err
	}

	actionSet := ComputeThemeActionSet(themeList, latestObjs)

	return types.ThemeActionSet{Items: actionSet}, nil
}

func RunThemeActionOnSet(executor execution.CommandExecutor, set []types.ThemeActionItem, actionToRun types.ThemeAction) error {
	for _, action := range set {
		// action to run is so we can run grouped actions per list
		if action.Action != actionToRun {
			continue
		}

		objectUrl, err := object_blueprint.GenerateStorageAccessUrl(action.Object)

		if err != nil {
			return errors.New(fmt.Sprintf("could not generate access hash for object id: %v, error: %s", action.Object.Id, err))
		}

		// execute action
		switch action.Action {
		case types.ThemeActionEnum.None:
			log.Debug("theme action none")
			break
		case types.ThemeActionEnum.Install:
			actionStr := fmt.Sprintf("wp theme install %s --activate --force --insecure", objectUrl)
			executor.ExecuteCommand([]string{actionStr})
			break
		case types.ThemeActionEnum.Upgrade:
			actionStr := fmt.Sprintf("wp theme install %s --activate --force --insecure", objectUrl)
			log.Infoln(actionStr)
			executor.ExecuteCommand([]string{actionStr})
			break
		case types.ThemeActionEnum.Downgrade:
			actionStr := fmt.Sprintf("wp theme install %s --activate --force --insecure", objectUrl)
			executor.ExecuteCommand([]string{actionStr})
			break
		case types.ThemeActionEnum.Uninstall:
			actionStr := fmt.Sprintf("wp theme uninstall %s", action.Name)
			executor.ExecuteCommand([]string{actionStr})
			break
		default:
			log.Fatal("default switch statement on theme sync command should not be reached")
		}

		time.Sleep(time.Millisecond * 250) // wait 250ms before each command
	}

	return nil
}

/*
 * This will iterate through the retrieved set and decide which action to take based on the desired set
 */
func ComputeThemeActionSet(themes []types.WpTheme, objects []types.ObjectBlueprint) []types.ThemeActionItem {

	sortSet := make(map[string]int)
	actionSet := make(map[int]types.ThemeActionItem)
	themeNameMap := make(map[string]types.ObjectBlueprint)

	// used to define the order of the themes according to the manifest file
	for i, object := range objects {
		sortSet[object.ExactName] = i
		themeNameMap[object.ExactName] = object
	}

	// add any other themes to this sort order
	sortSetLen := len(sortSet)
	for i, wp := range themes {
		if _, found := sortSet[wp.Name]; found == false {
			sortSet[wp.Name] = i + sortSetLen
		}
	}

	// all themes returned from server
	for _, theme := range themes {
		if dbTheme, exists := themeNameMap[theme.Name]; exists {

			// make a copy to prevent dbPlugin ref from being modified
			cpDbTheme := dbTheme

			dbSemver := semver.MustParse(cpDbTheme.Version)
			currentSemver := semver.MustParse(theme.Version)

			// the server version must match the manifest version
			if currentSemver.Equal(dbSemver) {
				actionSet[sortSet[theme.Name]] = types.ThemeActionItem{
					Name:   cpDbTheme.ExactName,
					Action: types.ThemeActionEnum.None,
					Object: &cpDbTheme,
				}
			} else if currentSemver.GreaterThan(dbSemver) {
				actionSet[sortSet[theme.Name]] = types.ThemeActionItem{
					Name:   cpDbTheme.ExactName,
					Action: types.ThemeActionEnum.Downgrade,
					Object: &cpDbTheme,
				}
			} else if currentSemver.LessThan(dbSemver) {
				actionSet[sortSet[theme.Name]] = types.ThemeActionItem{
					Name:   cpDbTheme.ExactName,
					Action: types.ThemeActionEnum.Upgrade,
					Object: &cpDbTheme,
				}
			}

		} else {

			// needs to be removed as it is not in the manifest
			actionSet[sortSet[theme.Name]] = types.ThemeActionItem{
				Name:   theme.Name,
				Action: types.ThemeActionEnum.Uninstall,
				Object: nil,
			}

		}
	}

	// loop through manifest themes, if they don't exist in the theme set then we need to install them
	for _, dbTheme := range objects {

		// make a copy to prevent dbPlugin ref from being modified
		cpDbTheme := dbTheme

		themeExists := false

		for _, sTheme := range themes {
			if sTheme.Name == cpDbTheme.ExactName {
				themeExists = true
				break
			}
		}

		// if theme exists is true, then we have already handled it in the above loop,
		// this loop is purely to handle the single case of the manifest theme existing only
		if themeExists == false {
			actionSet[sortSet[cpDbTheme.ExactName]] = types.ThemeActionItem{
				Name:   cpDbTheme.ExactName,
				Action: types.ThemeActionEnum.Install,
				Object: &cpDbTheme,
			}
		}
	}

	keys := make([]int, len(actionSet))
	i := 0
	for k := range actionSet {
		keys[i] = k
		i++
	}
	sort.Ints(keys)

	var sortedActionSet []types.ThemeActionItem

	for _, k := range keys {
		sortedActionSet = append(sortedActionSet, actionSet[k])
	}

	return sortedActionSet
}

func VerifyThemeActionExists(action string) {
	for _, item := range types.ThemeActionsList {
		if string(item) == action {
			return
		}
	}

	log.Fatal("Provided action does not exist: ", action)
}

func ParseThemeAction(action string) types.ThemeAction {
	switch action {
	case string(types.ThemeActionEnum.Install):
		return types.ThemeActionEnum.Install
	case string(types.ThemeActionEnum.Uninstall):
		return types.ThemeActionEnum.Uninstall
	case string(types.ThemeActionEnum.Upgrade):
		return types.ThemeActionEnum.Upgrade
	case string(types.ThemeActionEnum.Downgrade):
		return types.ThemeActionEnum.Downgrade
	case string(types.ThemeActionEnum.None):
		return types.ThemeActionEnum.None
	}

	return types.ThemeActionEnum.None
}

func ReverseThemeActionItemOrder(set []types.ThemeActionItem) []types.ThemeActionItem {
	size := reflect.ValueOf(set).Len()
	swap := reflect.Swapper(set)
	for i, j := 0, size-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}

	return set
}
