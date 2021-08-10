package types

import (
	"time"
)

type Account struct {
	Id        int64 `bun:"id,pk"`
	Uuid      string
	Name      string
	Key       string
	Enabled   bool
	Users     []*User   `bun:"m2m:users_accounts"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}
