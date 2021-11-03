package fixture

import "github.com/samvaughton/wpcommand/v2/pkg/db"

const TestAccountDefaultId = 1
const TestAccountAlphaId = 2
const TestAccountBetaId = 3
const TestAccountCharlieId = 4

func LoadTestAccounts() {
	accounts := []struct {
		Name string
		Key  string
	}{
		{
			Name: "Alpha",
			Key:  "alpha",
		},
		{
			Name: "Beta",
			Key:  "beta",
		},
		{
			Name: "Charlie",
			Key:  "charlie",
		},
	}

	for _, item := range accounts {
		_, err := db.AccountCreate(db.Db, item.Name, item.Key)

		if err != nil {
			panic(err)
		}
	}
}
