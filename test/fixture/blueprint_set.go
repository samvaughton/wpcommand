package fixture

import (
	"github.com/samvaughton/wpcommand/v2/pkg/db"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
)

const TestBlueprintSetAlphaId = 1
const TestBlueprintSetBetaId = 2
const TestBlueprintSetCharlieId = 3

func LoadTestBlueprintSets() {
	_, err := db.BlueprintCreateFromPayload(
		&types.CreateBlueprintSetPayload{
			Name: "Alpha",
		},
		TestAccountAlphaId,
	)

	if err != nil {
		panic(err)
	}

	_, err = db.BlueprintCreateFromPayload(
		&types.CreateBlueprintSetPayload{
			Name: "Beta",
		},
		TestAccountBetaId,
	)

	if err != nil {
		panic(err)
	}

	_, err = db.BlueprintCreateFromPayload(
		&types.CreateBlueprintSetPayload{
			Name: "Charlie",
		},
		TestAccountCharlieId,
	)

	if err != nil {
		panic(err)
	}
}
