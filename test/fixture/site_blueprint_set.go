package fixture

import "github.com/samvaughton/wpcommand/v2/pkg/db"

func AttachBlueprintSetsToSite() {
	err := db.SiteAddBlueprintSet(TestSiteAlpha1Id, TestBlueprintSetAlphaId)

	if err != nil {
		panic(err)
	}

	err = db.SiteAddBlueprintSet(TestSiteBetaId, TestBlueprintSetBetaId)

	if err != nil {
		panic(err)
	}

	err = db.SiteAddBlueprintSet(TestSiteCharlieId, TestBlueprintSetCharlieId)

	if err != nil {
		panic(err)
	}
}
