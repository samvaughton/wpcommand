package types

import (
	"time"
)

type User struct {
	Id        int64 `bun:"id,pk"`
	Uuid      string
	Email     string
	FirstName string
	LastName  string
	Password  string
	Enabled   bool
	Accounts  []*Account `bun:"m2m:users_accounts"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
