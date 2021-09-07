package types

import (
	"github.com/uptrace/bun"
	"time"
)

type ObjectBlueprintStorage struct {
	bun.BaseModel `bun:"object_blueprint_storage"`

	Id int64 `bun:"id,pk"`

	ObjectBlueprints []*ObjectBlueprint `bun:"m2m:object_blueprint_storage_relations"`

	Uuid string
	Hash string
	File []byte `bun:"type:bytea"`

	CreatedAt time.Time
}
