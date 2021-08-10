package types

import (
	"github.com/uptrace/bun"
	"time"
)

type BlueprintSet struct {
	Id        int64   `bun:"id,pk"`
	Sites     []*Site `bun:"m2m:sites_blueprint_sets"`
	Uuid      string
	Name      string
	Enabled   bool
	CreatedAt time.Time
	UpdatedAt bun.NullTime
}
