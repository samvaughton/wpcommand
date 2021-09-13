package types

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation"
	"io/ioutil"
	"net/http"
	"time"
)

type Account struct {
	Id           int64 `bun:"id,pk"`
	Uuid         string
	Name         string
	Key          string
	Enabled      bool
	UserAccounts []*UserAccount `bun:"rel:has-many"`
	CreatedAt    time.Time      `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt    time.Time      `bun:",nullzero,notnull,default:current_timestamp"`
}

type CreateAccountPayload struct {
	Name string
}

func NewCreateAccountPayloadFromHttpRequest(req *http.Request) (*CreateAccountPayload, error) {
	var item CreateAccountPayload

	bytes, err := ioutil.ReadAll(req.Body)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &item)

	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (p CreateAccountPayload) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Name, validation.Required),
	)
}

type UpdateAccountPayload struct {
	Name string
}

func NewUpdateAccountPayloadFromHttpRequest(req *http.Request) (*UpdateAccountPayload, error) {
	var item UpdateAccountPayload

	bytes, err := ioutil.ReadAll(req.Body)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &item)

	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (p UpdateAccountPayload) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Name, validation.Required),
	)
}
