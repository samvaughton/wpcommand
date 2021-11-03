package fixture

import (
	"github.com/samvaughton/wpcommand/v2/pkg/db"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
)

const TestSiteAlpha1Id = 1
const TestSiteAlpha2Id = 2
const TestSiteBetaId = 3
const TestSiteCharlieId = 4

func LoadTestSites() {

	sites := []types.Site{
		{
			Description:        "Alpha 1",
			AccountId:          TestAccountAlphaId,
			LabelSelector:      "k8.rentivo.com/name=alpha1",
			Namespace:          "test-alpha1",
			SiteEmail:          "test+alpha1@example.com",
			SiteUsername:       "alpha1",
			SitePassword:       "password",
			WpDomain:           "alpha1-wp.k8.rentivo.com",
			DockerRegistryName: "test-alpha1",
			Enabled:            true,
			TestMode:           true,
		},
		{
			Description:        "Alpha 2",
			AccountId:          TestAccountAlphaId,
			LabelSelector:      "k8.rentivo.com/name=alpha2",
			Namespace:          "test-alpha2",
			SiteEmail:          "test+alpha2@example.com",
			SiteUsername:       "alpha2",
			SitePassword:       "password",
			WpDomain:           "alpha2-wp.k8.rentivo.com",
			DockerRegistryName: "test-alpha2",
			Enabled:            true,
			TestMode:           true,
		},
		{
			Description:        "Beta",
			AccountId:          TestAccountBetaId,
			LabelSelector:      "k8.rentivo.com/name=beta",
			Namespace:          "test-beta",
			SiteEmail:          "test+beta@example.com",
			SiteUsername:       "beta",
			SitePassword:       "password",
			WpDomain:           "beta-wp.k8.rentivo.com",
			DockerRegistryName: "test-beta",
			Enabled:            true,
			TestMode:           true,
		},
		{
			Description:        "Charlie",
			AccountId:          TestAccountCharlieId,
			LabelSelector:      "k8.rentivo.com/name=charlie",
			Namespace:          "test-charlie",
			SiteEmail:          "test+charlie@example.com",
			SiteUsername:       "charlie",
			SitePassword:       "password",
			WpDomain:           "charlie-wp.k8.rentivo.com",
			DockerRegistryName: "test-charlie",
			Enabled:            true,
			TestMode:           true,
		},
	}

	for _, item := range sites {
		err := db.SiteCreateFromStruct(&item, item.AccountId)

		if err != nil {
			panic(err)
		}
	}
}
