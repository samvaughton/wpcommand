package types

import (
	"github.com/uptrace/bun"
)

type UserAccount struct {
	bun.BaseModel `bun:"users_accounts"`
	UserId        int64    `bun:"user_id,pk"`
	User          *User    `bun:"rel:has-one,join:user_id=id"`
	AccountId     int64    `bun:"account_id,pk"`
	Account       *Account `bun:"rel:has-one,join:account_id=id"`
}
