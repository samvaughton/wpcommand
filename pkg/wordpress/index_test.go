package wordpress

import (
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

const lazyblocksOutput = `[
	{"ID":292880,"post_title":"Featurette","post_name":"featurette","post_date":"2021-08-11 14:29:44","post_status":"publish"},
	{"ID":292881,"post_title":"Row Block","post_name":"row-block","post_date":"2021-08-11 14:29:44","post_status":"publish"},
	{"ID":292882,"post_title":"Card","post_name":"card","post_date":"2021-08-11 14:29:44","post_status":"publish"}
]`

const acfOutput = `[
	{"ID":292906,"post_title":"Options: Redirects","post_name":"group_5f439fd8b9389","post_date":"2021-08-11 14:29:46","post_status":"publish"},
	{"ID":292908,"post_title":"Options: Translations","post_name":"group_5ee21d4e48f52","post_date":"2021-08-11 14:29:46","post_status":"publish"},
	{"ID":292910,"post_title":"Page","post_name":"group_5e3d2ca89944c","post_date":"2021-08-11 14:29:46","post_status":"publish"}
]`

const postAndPagesOutput = `[
	{"ID":286632,"post_title":"4 Reasons to Add the Canadian Rockies to Your Bucket List","post_name":"4-reasons-visit-canadian-rockies","post_date":"2019-10-18 16:09:13","post_status":"publish"}
]`

func TestGetSiteLazyblocksPosts(t *testing.T) {
	expected := []types.WpPost{
		types.WpPost{
			Id:     292880,
			Title:  "Featurette",
			Name:   "featurette",
			Date:   "2021-08-11 14:29:44",
			Status: "publish",
		},
		types.WpPost{
			Id:     292881,
			Title:  "Row Block",
			Name:   "row-block",
			Date:   "2021-08-11 14:29:44",
			Status: "publish",
		},
		types.WpPost{
			Id:     292882,
			Title:  "Card",
			Name:   "card",
			Date:   "2021-08-11 14:29:44",
			Status: "publish",
		},
	}

	actual, err := GetSiteLazyblocksPosts(newTestExecutor(lazyblocksOutput))

	t.Run("TestGetSiteLazyblocksPosts", func(t *testing.T) {
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestGetSitePostAndPages(t *testing.T) {
	expected := []types.WpPost{
		types.WpPost{
			Id:     286632,
			Title:  "4 Reasons to Add the Canadian Rockies to Your Bucket List",
			Name:   "4-reasons-visit-canadian-rockies",
			Date:   "2019-10-18 16:09:13",
			Status: "publish",
		},
	}

	actual, err := GetSitePostAndPages(newTestExecutor(postAndPagesOutput))

	t.Run("TestGetSitePostAndPages", func(t *testing.T) {
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestGetSiteAcfFieldGroups(t *testing.T) {
	expected := []types.WpPost{
		types.WpPost{
			Id:     292906,
			Title:  "Options: Redirects",
			Name:   "group_5f439fd8b9389",
			Date:   "2021-08-11 14:29:46",
			Status: "publish",
		},
		types.WpPost{
			Id:     292908,
			Title:  "Options: Translations",
			Name:   "group_5ee21d4e48f52",
			Date:   "2021-08-11 14:29:46",
			Status: "publish",
		},
		types.WpPost{
			Id:     292910,
			Title:  "Page",
			Name:   "group_5e3d2ca89944c",
			Date:   "2021-08-11 14:29:46",
			Status: "publish",
		},
	}

	actual, err := GetSiteAcfFieldGroups(newTestExecutor(acfOutput))

	t.Run("TestGetSiteAcfFieldGroups", func(t *testing.T) {
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
}
