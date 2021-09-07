package db

import (
	"context"
	"github.com/google/uuid"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	"time"
)

func BlueprintObjectStorageGetByHash(hash string) (*types.ObjectBlueprintStorage, error) {
	item := new(types.ObjectBlueprintStorage)

	err := Db.
		NewSelect().
		Model(item).
		Relation("ObjectBlueprints").
		Where("\"object_blueprint_storage\".hash = ?", hash).
		Scan(context.Background())

	if err != nil {
		return nil, err
	}

	return item, nil
}

func BlueprintObjectStorageCreateFromBytes(hash string, data []byte) (*types.ObjectBlueprintStorage, error) {
	ob := &types.ObjectBlueprintStorage{
		Uuid:      uuid.New().String(),
		CreatedAt: time.Now(),
		Hash:      hash,
		File:      data,
	}
	_, err := Db.NewInsert().Model(ob).Returning("*").Exec(context.Background())

	return ob, err
}
