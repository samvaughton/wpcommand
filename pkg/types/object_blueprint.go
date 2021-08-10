package types

import (
	"github.com/uptrace/bun"
	"time"
)

const ObjectBlueprintTypePlugin = "PLUGIN"
const ObjectBlueprintTypeTheme = "THEME"

type ObjectBlueprint struct {
	Id         int64 `bun:"id,pk"`
	RevisionId int64

	BlueprintSetId int64
	BlueprintSet   *BlueprintSet `bun:"rel:belongs-to"`

	Uuid     string
	SetOrder int
	Type     string
	Name     string
	Enabled  bool

	Version   string
	ExactName string

	OriginalObjectUrl   string
	VersionedObjectUrl  string
	VersionedObjectHash string

	CreatedAt time.Time
	UpdatedAt bun.NullTime
}
