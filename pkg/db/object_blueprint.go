package db

import (
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
	"time"
)

func BlueprintObjectGetAllRevisionsSafe(objectUuid string, blueprintId int64) ([]*types.ObjectBlueprint, error) {
	var err error
	var items = make([]*types.ObjectBlueprint, 0)

	err = Db.
		NewSelect().
		Model(&items).
		Relation("BlueprintSet").
		Where("\"object_blueprint\".uuid = ? AND blueprint_set_id = ?", objectUuid, blueprintId).
		Order("revision_id DESC").
		Scan(context.Background())

	if err != nil {
		log.Error(err)

		return items, err
	}

	return items, nil
}

func BlueprintObjectGetByUuidAndRevisionSafe(objectUuid string, revisionId int, blueprintId int64) (*types.ObjectBlueprint, error) {
	item := new(types.ObjectBlueprint)

	err := Db.
		NewSelect().
		Model(item).
		Relation("BlueprintSet").
		Relation("ObjectBlueprintStorage").
		Where("\"object_blueprint\".uuid = ? and revision_id = ? and blueprint_set_id = ?", objectUuid, revisionId, blueprintId).
		Scan(context.Background())

	if err != nil {
		return nil, err
	}

	return item, nil
}

func BlueprintObjectGetLatestRevisionByUuid(objectUuid string) (*types.ObjectBlueprint, error) {
	item := new(types.ObjectBlueprint)

	err := Db.
		NewSelect().
		Model(item).
		Relation("BlueprintSet").
		Where(
			"\"object_blueprint\".uuid = ? "+
				"AND revision_id = (SELECT max(revision_id) FROM object_blueprints AS ob2 WHERE \"object_blueprint\".uuid = ob2.uuid)",
			objectUuid,
		).
		Limit(1).
		Scan(context.Background())

	if err != nil {
		return nil, err
	}

	return item, nil
}

func BlueprintObjectGetLatestRevisionByUuidSafe(objectUuid string, blueprintId int64) (*types.ObjectBlueprint, error) {
	item, err := BlueprintObjectGetLatestRevisionByUuid(objectUuid)

	if err != nil {
		return nil, err
	}

	if item.BlueprintSetId != blueprintId {
		return nil, errors.New("could not locate latest blueprint object")
	}

	return item, nil
}

func BlueprintObjectGetByUuidAndRevision(objectUuid string, revisionId int) (*types.ObjectBlueprint, error) {
	item := new(types.ObjectBlueprint)

	err := Db.NewSelect().Model(item).Where("uuid = ? and revision_id = ?", objectUuid, revisionId).Scan(context.Background())

	if err != nil {
		return nil, err
	}

	return item, nil
}

func GetLatestObjectBlueprintsForSiteAndType(siteId int64, bpType string) []types.ObjectBlueprint {
	var err error
	var items []types.ObjectBlueprint

	err = Db.
		NewSelect().
		Model(&items).
		Relation("ObjectBlueprintStorage").
		Relation("BlueprintSet").
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

func BlueprintObjectCreateFromPayload(tx bun.IDB, payload *types.CreateObjectBlueprintPayload, blueprintSetId int64) (*types.ObjectBlueprint, error) {
	ob := &types.ObjectBlueprint{
		Uuid:              uuid.New().String(),
		CreatedAt:         time.Now(),
		UpdatedAt:         bun.NullTime{Time: time.Now()},
		BlueprintSetId:    blueprintSetId,
		RevisionId:        0,
		Active:            true, // new object first revision set to true
		Type:              payload.Type,
		Name:              payload.Name,
		ExactName:         payload.ExactName,
		Version:           payload.Version,
		OriginalObjectUrl: payload.Url,
		SetOrder:          payload.SetOrder,
	}

	_, err := tx.NewInsert().Model(ob).Returning("*").Exec(context.Background())

	return ob, err
}

func BlueprintObjectCreateNewRevisionFromPayload(tx bun.IDB, objectUuid string, payload *types.UpdatedVersionObjectBlueprintPayload) (*types.ObjectBlueprint, error) {
	// we need to get the latest object for this uuid
	latest, err := BlueprintObjectGetLatestRevisionByUuid(objectUuid)

	if err != nil {
		return nil, err
	}

	url := payload.Url

	if len(url) == 0 {
		url = latest.OriginalObjectUrl
	}

	ob := &types.ObjectBlueprint{
		Uuid:              latest.Uuid,
		CreatedAt:         time.Now(),
		UpdatedAt:         bun.NullTime{Time: time.Now()},
		BlueprintSetId:    latest.BlueprintSetId,
		RevisionId:        latest.RevisionId + 1,
		Active:            true,
		Type:              latest.Type,
		Name:              latest.Name,
		ExactName:         latest.ExactName,
		Version:           payload.Version,
		OriginalObjectUrl: url,
	}

	_, err = tx.NewInsert().Model(ob).Returning("*").Exec(context.Background())

	return ob, err
}
