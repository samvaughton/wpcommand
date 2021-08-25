package db

import (
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	"github.com/uptrace/bun"
	"time"
)

func BlueprintGetByUuid(uuid string) (*types.BlueprintSet, error) {
	item := new(types.BlueprintSet)

	err := Db.NewSelect().Model(item).Where("uuid = ?", uuid).Scan(context.Background())

	if err != nil {
		return nil, err
	}

	return item, nil
}

func BlueprintGetByUuidSafe(uuid string, accountId int64) (*types.BlueprintSet, error) {
	item, err := BlueprintGetByUuid(uuid)

	if err != nil {
		return nil, err
	}

	if item.AccountId != accountId {
		return nil, errors.New("not found")
	}

	return item, nil
}

func BlueprintsGetForAccountId(accountId int64) ([]*types.BlueprintSet, error) {
	items := make([]*types.BlueprintSet, 0)

	err := Db.
		NewSelect().
		Model(&items).
		Relation("Account").
		Where("account_id = ? ", accountId).
		Order("created_at DESC").
		Scan(context.Background())

	if err != nil {
		return nil, err
	}

	return items, nil
}

func BlueprintsGetForSiteSafe(siteId int64, accountId int64) ([]*types.BlueprintSet, error) {
	items := make([]*types.BlueprintSet, 0)

	err := Db.
		NewSelect().
		Model(&items).
		Relation("Sites").
		Join("JOIN sites_blueprint_sets AS sbs ON \"blueprint_set\".\"id\" = sbs.site_id").
		Where("account_id = ? and sbs.site_id = ?", accountId, siteId).
		Order("name ASC").
		Scan(context.Background())

	if err != nil {
		return nil, err
	}

	return items, nil
}

func BlueprintCreateFromPayload(payload *types.CreateBlueprintSetPayload, accountId int64) (*types.BlueprintSet, error) {
	return BlueprintCreate(&types.BlueprintSet{
		Uuid:      uuid.New().String(),
		Name:      payload.Name,
		Enabled:   true,
		CreatedAt: time.Now(),
		UpdatedAt: bun.NullTime{Time: time.Now()},
	}, accountId)
}

func BlueprintCreate(bp *types.BlueprintSet, accountId int64) (*types.BlueprintSet, error) {
	bp.AccountId = accountId

	_, err := Db.NewInsert().Model(bp).Returning("*").Exec(context.Background())

	if err != nil {
		return nil, err
	}

	return bp, nil
}
