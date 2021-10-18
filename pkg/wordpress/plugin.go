package wordpress

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"
	"github.com/samvaughton/wpcommand/v2/pkg/execution"
	"github.com/samvaughton/wpcommand/v2/pkg/object_blueprint"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
	"reflect"
	"sort"
	"strings"
	"time"
)

func GetSitePluginStatuses(executor execution.CommandExecutor, latestObjs []types.ObjectBlueprint) (types.PluginActionSet, error) {
	result, err := executor.ExecuteCommand([]string{"wp plugin list --format=json"})

	if err != nil {
		return types.PluginActionSet{}, err
	}

	pluginList, err := PluginListFromJson(result.Output)

	if err != nil {
		return types.PluginActionSet{}, err
	}

	actionSet := ComputePluginActionSet(pluginList, latestObjs)

	return types.PluginActionSet{Items: actionSet}, nil
}

func RunPluginActionOnSet(executor execution.CommandExecutor, set []types.PluginActionItem, actionToRun types.PluginAction) error {
	var errs []string
	for _, action := range set {
		// action to run is so we can run grouped actions per list
		if action.Action != actionToRun {
			continue
		}

		objectUrl, err := object_blueprint.GenerateStorageAccessUrl(action.Object)

		if err != nil {
			return errors.New(fmt.Sprintf("could not generate file access hash url for plugin: %v, error: %s", action.Name, err))
		}

		err = nil // clear

		// execute action
		switch action.Action {
		case types.PluginActionEnum.None:
			log.Debug("plugin action none")
			break
		case types.PluginActionEnum.Install:
			actionStr := fmt.Sprintf("wp plugin install %s --activate --force --insecure", objectUrl)
			_, err = executor.ExecuteCommand([]string{actionStr})
			break
		case types.PluginActionEnum.Upgrade:
			actionStr := fmt.Sprintf("wp plugin install %s --activate --force --insecure", objectUrl)
			log.Infoln(actionStr)
			_, err = executor.ExecuteCommand([]string{actionStr})
			break
		case types.PluginActionEnum.Downgrade:
			actionStr := fmt.Sprintf("wp plugin install %s --activate --force --insecure", objectUrl)
			_, err = executor.ExecuteCommand([]string{actionStr})
			break
		case types.PluginActionEnum.Uninstall:
			actionStr := fmt.Sprintf("wp plugin uninstall --deactivate %s", action.Name)
			_, err = executor.ExecuteCommand([]string{actionStr})
			break
		default:
			log.Fatal("default switch statement on plugin sync command should not be reached")
		}

		if err != nil {
			errs = append(errs, err.Error())
		}

		time.Sleep(time.Millisecond * 250) // wait 250ms before each command
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ", "))
	}

	return nil
}

/*
 * This will iterate through the retrieved set and decide which action to take based on the desired set
 */
func ComputePluginActionSet(plugins []types.WpPlugin, objects []types.ObjectBlueprint) []types.PluginActionItem {

	sortSet := make(map[string]int)
	actionSet := make(map[int]types.PluginActionItem)
	pluginNameMap := make(map[string]types.ObjectBlueprint)

	// used to define the order of the plugins according to the manifest file
	for i, object := range objects {
		sortSet[object.ExactName] = i
		pluginNameMap[object.ExactName] = object
	}

	// add any other plugins to this sort order
	sortSetLen := len(sortSet)
	for i, wp := range plugins {
		if _, found := sortSet[wp.Name]; found == false {
			sortSet[wp.Name] = i + sortSetLen
		}
	}

	// all plugins returned from server
	for _, plugin := range plugins {
		if dbPlugin, exists := pluginNameMap[plugin.Name]; exists {

			cpDbPlugin := dbPlugin

			dbSemver := semver.MustParse(cpDbPlugin.Version)
			currentSemver := semver.MustParse(plugin.Version)

			if currentSemver.Equal(dbSemver) {
				actionSet[sortSet[plugin.Name]] = types.PluginActionItem{
					Name:   cpDbPlugin.ExactName,
					Action: types.PluginActionEnum.None,
					Object: &cpDbPlugin,
				}
			} else if currentSemver.GreaterThan(dbSemver) {
				actionSet[sortSet[plugin.Name]] = types.PluginActionItem{
					Name:   cpDbPlugin.ExactName,
					Action: types.PluginActionEnum.Downgrade,
					Object: &cpDbPlugin,
				}
			} else if currentSemver.LessThan(dbSemver) {
				actionSet[sortSet[plugin.Name]] = types.PluginActionItem{
					Name:   cpDbPlugin.ExactName,
					Action: types.PluginActionEnum.Upgrade,
					Object: &cpDbPlugin,
				}
			}

		} else {

			// needs to be removed as it is not in the manifest
			actionSet[sortSet[plugin.Name]] = types.PluginActionItem{
				Name:   plugin.Name,
				Action: types.PluginActionEnum.Uninstall,
				Object: nil,
			}

		}
	}

	// loop through manifest plugins, if they don't exist in the plugin set then we need to install them
	for _, dbPlugin := range objects {

		// make a copy to prevent dbPlugin ref from being modified
		cpDbPlugin := dbPlugin

		pluginExists := false

		for _, sPlugin := range plugins {
			if sPlugin.Name == cpDbPlugin.ExactName {
				pluginExists = true
				break
			}
		}

		// if plugin exists is true, then we have already handled it in the above loop,
		// this loop is purely to handle the single case of the manifest plugin existing only
		if pluginExists == false {

			actionSet[sortSet[dbPlugin.ExactName]] = types.PluginActionItem{
				Name:   cpDbPlugin.ExactName,
				Action: types.PluginActionEnum.Install,
				Object: &cpDbPlugin,
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

	var sortedActionSet []types.PluginActionItem

	for _, k := range keys {
		sortedActionSet = append(sortedActionSet, actionSet[k])
	}

	return sortedActionSet
}

func VerifyPluginActionExists(action string) {
	for _, item := range types.PluginActionsList {
		if string(item) == action {
			return
		}
	}

	log.Fatal("Provided action does not exist: ", action)
}

func ParsePluginAction(action string) types.PluginAction {
	switch action {
	case string(types.PluginActionEnum.Install):
		return types.PluginActionEnum.Install
	case string(types.PluginActionEnum.Uninstall):
		return types.PluginActionEnum.Uninstall
	case string(types.PluginActionEnum.Upgrade):
		return types.PluginActionEnum.Upgrade
	case string(types.PluginActionEnum.Downgrade):
		return types.PluginActionEnum.Downgrade
	case string(types.PluginActionEnum.None):
		return types.PluginActionEnum.None
	}

	return types.PluginActionEnum.None
}

func ReversePluginActionItemOrder(set []types.PluginActionItem) []types.PluginActionItem {
	size := reflect.ValueOf(set).Len()
	swap := reflect.Swapper(set)
	for i, j := 0, size-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}

	return set
}
