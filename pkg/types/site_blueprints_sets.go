package types

import (
	"github.com/uptrace/bun"
)

type SiteBlueprintSet struct {
	bun.BaseModel  `bun:"sites_blueprint_sets"`
	SiteId         int64         `bun:"site_id,pk"`
	Site           *Site         `bun:"rel:has-one,join:site_id=id"`
	BlueprintSetId int64         `bun:"blueprint_set_id,pk"`
	BlueprintSet   *BlueprintSet `bun:"rel:has-one,join:blueprint_set_id=id"`
}
