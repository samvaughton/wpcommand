//go:build integration
// +build integration

package registry

import (
	"encoding/json"
	"github.com/samvaughton/wpcommand/v2/pkg/config"
	"github.com/samvaughton/wpcommand/v2/pkg/db"
	"github.com/samvaughton/wpcommand/v2/pkg/execution"
	"github.com/samvaughton/wpcommand/v2/pkg/pipeline"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	"github.com/samvaughton/wpcommand/v2/test/fixture"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPluginCommands(t *testing.T) {
	config.InitConfigFromFile("../../config.test.yaml")
	db.InitDbConnection()
	defer db.Db.Close()

	testSite := db.SiteMustGetById(fixture.TestSiteAlpha1Id)

	initDebugExecutor := func(site *types.Site, output string) *execution.DebugCommandExecutor {
		exec, err := execution.NewCommandExecutor(testSite, &types.Config{})
		assert.Nil(t, err)
		debugExec, valid := exec.(*execution.DebugCommandExecutor)
		assert.True(t, valid)

		debugExec.MockOutput = output

		return debugExec
	}

	t.Run("TestPluginStatusCommandNoneInstalled", func(t *testing.T) {
		exec := initDebugExecutor(testSite, `[]`)

		defer exec.ClearLog()

		cmds := []pipeline.SiteCommand{
			GetPluginsStatusCommand(testSite),
		}

		pl := &pipeline.SiteCommandPipeline{
			Name:     "test pipeline",
			Site:     testSite,
			Executor: exec,
			Options:  pipeline.ExecuteOptions{},
			Commands: cmds,
		}

		pl.Run()

		assert.Empty(t, pl.Errors)
		assert.Contains(t, exec.CommandLog, "wp plugin list --format=json")

		statusResult := pl.Results[0]
		actionSet, valid := statusResult.Data.(types.PluginActionSet)

		assert.True(t, valid)

		expected := `[{"Name":"rentivo-advanced-custom-fields-pro","Action":"install"},{"Name":"acf-code-field","Action":"install"},{"Name":"acf-to-rest-api","Action":"install"},{"Name":"acf-map-bbox","Action":"install"},{"Name":"polylang-pro","Action":"install"},{"Name":"wp-graphql","Action":"install"},{"Name":"wp-graphql-polylang-0.5.0","Action":"install"},{"Name":"wp-graphql-gutenberg-develop","Action":"install"},{"Name":"wp-graphql-acf-0.4.1","Action":"install"},{"Name":"wp-graphql-meta-query-0.1.0","Action":"install"},{"Name":"wp-graphiql-master","Action":"install"},{"Name":"wp-gatsby","Action":"install"},{"Name":"lazy-blocks","Action":"install"},{"Name":"svg-support","Action":"install"},{"Name":"duplicate-page","Action":"install"},{"Name":"wordpress-importer","Action":"install"},{"Name":"rentivo-simba","Action":"install"}]`

		jsonOutput, _ := json.Marshal(actionSet.Items)
		assert.JSONEq(t, expected, string(jsonOutput))
	})

}
