package types

import (
	"github.com/uptrace/bun"
)

type ObjectBlueprintStorageRelation struct {
	bun.BaseModel            `bun:"object_blueprint_storage_relations"`
	ObjectBlueprintId        int64                   `bun:"object_blueprint_id,pk"`
	ObjectBlueprint          *ObjectBlueprint        `bun:"rel:has-one,join:object_blueprint_id=id"`
	ObjectBlueprintStorageId int64                   `bun:"object_blueprint_storage_id,pk"`
	ObjectBlueprintStorage   *ObjectBlueprintStorage `bun:"rel:has-one,join:object_blueprint_storage_id=id"`
}
