package types

import (
	"fmt"
	"github.com/uptrace/bun"
)

type UserAccount struct {
	bun.BaseModel `bun:"users_accounts"`
	Uuid          string   `bun:",unique"`
	Roles         []string `bun:",array"`
	UserId        int64    `bun:"user_id,pk"`
	User          *User    `bun:"rel:has-one,join:user_id=id"`
	AccountId     int64    `bun:"account_id,pk"`
	Account       *Account `bun:"rel:has-one,join:account_id=id"`
}

func (ua *UserAccount) GetCasbinPolicyKey() string {
	return fmt.Sprintf("u%v_a%v", ua.UserId, ua.AccountId)
}
