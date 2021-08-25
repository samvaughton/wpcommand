package types

import (
	"gopkg.in/guregu/null.v3"
	"time"
)

const CommandTypeHttpCall = "HTTP_CALL"
const CommandTypeWpBuiltIn = "WP_BUILT_IN"

type Command struct {
	Id          int64 `bun:"id,pk"`
	AccountId   null.Int
	Account     *Account `bun:"rel:belongs-to" json:"-"`
	SiteId      null.Int
	Site        *Site `bun:"rel:belongs-to" json:"-"`
	Public      bool
	Uuid        string
	Key         string
	Type        string
	Description string
	HttpMethod  string
	HttpUrl     string
	HttpHeaders string
	HttpBody    string
	CreatedAt   time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}

func (c *Command) IsDefault() bool {
	return c.AccountId.Valid == false && c.SiteId.Valid == false
}
