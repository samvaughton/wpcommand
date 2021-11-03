package main

import (
	"github.com/samvaughton/wpcommand/v2/pkg/auth"
	"github.com/samvaughton/wpcommand/v2/pkg/config"
	"github.com/samvaughton/wpcommand/v2/pkg/db"
	"github.com/samvaughton/wpcommand/v2/pkg/flow"
	"github.com/samvaughton/wpcommand/v2/pkg/registry"
	"github.com/samvaughton/wpcommand/v2/pkg/util"
	"github.com/samvaughton/wpcommand/v2/test/fixture"
	"os"
)

/*
 * Although we could load various fixtures via a SQL file, initializing it this way
 * allows us to ensure our db insert functions also work as intended
 */

func main() {
	configData, err := os.ReadFile("./config.test.yaml")

	if err != nil {
		panic(err)
	}

	authData, err := os.ReadFile("./casbin/model.conf")

	if err != nil {
		panic(err)
	}

	config.InitConfig(string(configData))
	db.InitDbConnection()
	util.SetupLogging()
	auth.InitAuth(string(authData))
	flow.CreateDefaultAccountAndUser()
	registry.CreateDefaultCommands()

	fixture.LoadTestAccounts()
	fixture.LoadTestSites()
	fixture.LoadTestBlueprintSets()
	fixture.LoadTestObjectBlueprints()

	fixture.AttachBlueprintSetsToSite()
}
