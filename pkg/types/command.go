package types

import (
	"time"
)

type Command struct {
	Id          int64 `bun:"id,pk"`
	AccountId   int64
	Account     *Account `bun:"rel:belongs-to"`
	Uuid        string
	Key         string
	Description string
	HttpMethod  string
	HttpUrl     string
	HttpHeaders string
	HttpBody    string
	CreatedAt   time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}
