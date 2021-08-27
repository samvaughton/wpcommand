package db

import (
	"context"
	"github.com/google/uuid"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
	"time"
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

func GetLatestBlueprintObjectsForBlueprintSetId(blueprintSetId int64) ([]*types.ObjectBlueprint, error) {
	var err error
	var items = make([]*types.ObjectBlueprint, 0)

	err = Db.
		NewSelect().
		Model(&items).
		Join("JOIN blueprint_sets AS bs ON \"object_blueprint\".blueprint_set_id = bs.id").
		Where("\"object_blueprint\".blueprint_set_id = ?", blueprintSetId).
		Where("revision_id = (SELECT max(revision_id) FROM object_blueprints AS ob2 WHERE \"object_blueprint\".uuid = ob2.uuid)").
		Order("set_order ASC").
		Scan(context.Background())

	if err != nil {
		log.Error(err)
	}

	return items, nil
}

func BlueprintObjectCreateFromPayload(payload *types.CreateObjectBlueprintPayload, blueprintSetId int64) (*types.ObjectBlueprint, error) {
	ob := &types.ObjectBlueprint{
		Uuid:              uuid.New().String(),
		CreatedAt:         time.Now(),
		UpdatedAt:         bun.NullTime{Time: time.Now()},
		BlueprintSetId:    blueprintSetId,
		RevisionId:        0,
		Enabled:           true,
		Name:              payload.Name,
		ExactName:         payload.ExactName,
		Version:           payload.Version,
		OriginalObjectUrl: payload.Url,
	}
	_, err := Db.NewInsert().Model(&ob).Returning("*").Exec(context.Background())

	return ob, err
}
