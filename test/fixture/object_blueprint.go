package fixture

import (
	"github.com/samvaughton/wpcommand/v2/pkg/db"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
)

func LoadTestObjectBlueprints() {

	items := []types.CreateObjectBlueprintPayload{
		{
			Type:      types.ObjectBlueprintTypePlugin,
			Name:      "advanced-custom-fields-pro",
			ExactName: "rentivo-advanced-custom-fields-pro",
			Version:   "5.9.5",
			SetOrder:  10,
			Url:       "https://github.com/aptenex/k8-store/raw/main/plugins/advanced-custom-fields-pro.zip",
		},
		{
			Type:      types.ObjectBlueprintTypePlugin,
			Name:      "acf-code-field",
			ExactName: "acf-code-field",
			Version:   "1.8",
			SetOrder:  20,
			Url:       "https://github.com/aptenex/k8-store/raw/main/plugins/acf-code-field.zip",
		},
		{
			Type:      types.ObjectBlueprintTypePlugin,
			Name:      "acf-to-rest-api",
			ExactName: "acf-to-rest-api",
			Version:   "3.3.2",
			SetOrder:  30,
			Url:       "https://github.com/aptenex/k8-store/raw/main/plugins/acf-to-rest-api.zip",
		},
		{
			Type:      types.ObjectBlueprintTypePlugin,
			Name:      "acf-map-bbox",
			ExactName: "acf-map-bbox",
			Version:   "1.0.0",
			SetOrder:  40,
			Url:       "https://github.com/aptenex/k8-store/raw/main/plugins/acf-map-bbox.zip",
		},
		{
			Type:      types.ObjectBlueprintTypePlugin,
			Name:      "polylang-pro",
			ExactName: "polylang-pro",
			Version:   "3.0.3",
			SetOrder:  50,
			Url:       "https://github.com/aptenex/k8-store/raw/main/plugins/polylang-pro.zip",
		},
		{
			Type:      types.ObjectBlueprintTypePlugin,
			Name:      "wp-graphql",
			ExactName: "wp-graphql",
			Version:   "1.3.9",
			SetOrder:  60,
			Url:       "https://github.com/aptenex/k8-store/raw/main/plugins/wp-graphql.zip",
		},
		{
			Type:      types.ObjectBlueprintTypePlugin,
			Name:      "wp-graphql-polylang-0.5.0",
			ExactName: "wp-graphql-polylang-0.5.0",
			Version:   "0.5.0",
			SetOrder:  70,
			Url:       "https://github.com/aptenex/k8-store/raw/main/plugins/wp-graphql-polylang-0.5.0.zip",
		},
		{
			Type:      types.ObjectBlueprintTypePlugin,
			Name:      "wp-graphql-gutenberg-develop",
			ExactName: "wp-graphql-gutenberg-develop",
			Version:   "0.3.8",
			SetOrder:  80,
			Url:       "https://github.com/aptenex/k8-store/raw/main/plugins/wp-graphql-gutenberg-develop.zip",
		},
		{
			Type:      types.ObjectBlueprintTypePlugin,
			Name:      "wp-graphql-acf-0.4.1",
			ExactName: "wp-graphql-acf-0.4.1",
			Version:   "0.4.1",
			SetOrder:  90,
			Url:       "https://github.com/aptenex/k8-store/raw/main/plugins/wp-graphql-acf-0.4.1.zip",
		},
		{
			Type:      types.ObjectBlueprintTypePlugin,
			Name:      "wp-graphql-meta-query-0.1.0",
			ExactName: "wp-graphql-meta-query-0.1.0",
			Version:   "0.1.0",
			SetOrder:  100,
			Url:       "https://github.com/aptenex/k8-store/raw/main/plugins/wp-graphql-meta-query-0.1.0.zip",
		},
		{
			Type:      types.ObjectBlueprintTypePlugin,
			Name:      "wp-graphiql-master",
			ExactName: "wp-graphiql-master",
			Version:   "1.0.1",
			SetOrder:  110,
			Url:       "https://github.com/aptenex/k8-store/raw/main/plugins/wp-graphiql-master.zip",
		},
		{
			Type:      types.ObjectBlueprintTypePlugin,
			Name:      "wp-gatsby",
			ExactName: "wp-gatsby",
			Version:   "1.0.10",
			SetOrder:  120,
			Url:       "https://github.com/aptenex/k8-store/raw/main/plugins/wp-gatsby.zip",
		},
		{
			Type:      types.ObjectBlueprintTypePlugin,
			Name:      "lazy-blocks",
			ExactName: "lazy-blocks",
			Version:   "2.3.1",
			SetOrder:  130,
			Url:       "https://github.com/aptenex/k8-store/raw/main/plugins/lazy-blocks.zip",
		},
		{
			Type:      types.ObjectBlueprintTypePlugin,
			Name:      "svg-support",
			ExactName: "svg-support",
			Version:   "2.3.18",
			SetOrder:  140,
			Url:       "https://github.com/aptenex/k8-store/raw/main/plugins/svg-support.zip",
		},
		{
			Type:      types.ObjectBlueprintTypePlugin,
			Name:      "duplicate-page",
			ExactName: "duplicate-page",
			Version:   "4.4",
			SetOrder:  150,
			Url:       "https://github.com/aptenex/k8-store/raw/main/plugins/duplicate-page.zip",
		},
		{
			Type:      types.ObjectBlueprintTypePlugin,
			Name:      "wordpress-importer",
			ExactName: "wordpress-importer",
			Version:   "0.7",
			SetOrder:  160,
			Url:       "https://github.com/aptenex/k8-store/raw/main/plugins/wordpress-importer.zip",
		},
		{
			Type:      types.ObjectBlueprintTypePlugin,
			Name:      "Rentivo Simba",
			ExactName: "rentivo-simba",
			Version:   "1.0.8",
			SetOrder:  170,
			Url:       "https://github.com/aptenex/k8-store/raw/main/plugins/rentivo-simba.zip",
		},
	}

	createForBlueprintSetId(items, TestBlueprintSetAlphaId)
	createForBlueprintSetId(items, TestBlueprintSetBetaId)
	createForBlueprintSetId(items, TestBlueprintSetCharlieId)
}

func createForBlueprintSetId(items []types.CreateObjectBlueprintPayload, blueprintSetId int64) {
	for _, item := range items {
		_, err := db.BlueprintObjectCreateFromPayload(db.Db, &item, blueprintSetId)

		if err != nil {
			panic(err)
		}
	}
}
