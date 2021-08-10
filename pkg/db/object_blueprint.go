package db

import (
	"context"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
)

func GetLatestObjectBlueprintsForSiteAndType(siteId int64, bpType string) []types.ObjectBlueprint {
	var err error
	var items []types.ObjectBlueprint

	err = Db.
		NewSelect().
		Model(&items).
		Join("JOIN blueprint_sets AS bs ON \"object_blueprint\".blueprint_set_id = bs.id").
		Join("JOIN sites_blueprint_sets AS sbs ON bs.id = sbs.blueprint_set_id").
		Where("\"object_blueprint\".type = ?", bpType).
		Where("sbs.site_id = ?", siteId).
		// not the most efficient query
		Where("revision_id = (SELECT max(revision_id) FROM object_blueprints AS ob2 WHERE \"object_blueprint\".uuid = ob2.uuid)").
		Order("set_order ASC").
		Scan(context.Background())

	if err != nil {
		log.Error(err)
	}

	return items
}
