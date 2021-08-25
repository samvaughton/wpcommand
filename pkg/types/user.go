package types

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"io/ioutil"
	"net/http"
	"time"
)

type User struct {
	Id           int64 `bun:"id,pk" json:"-"`
	Uuid         string
	Email        string
	FirstName    string
	LastName     string
	Password     string `json:"-"`
	Enabled      bool
	SuperAdmin   bool           `json:"-"`
	UserAccounts []*UserAccount `bun:"rel:has-many" json:"-"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserCreatePayload struct {
	Email     string `validate:"required,email"`
	FirstName string `validate:"required"`
	LastName  string `validate:"required"`
	Password  string `validate:"required,gt=8"`
}

func (p UserCreatePayload) Validate(uniqueCheck validation.Rule) error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Email, validation.Required, is.Email, uniqueCheck),
		validation.Field(&p.FirstName, validation.Required),
		validation.Field(&p.LastName, validation.Required),
		validation.Field(&p.Password, validation.Required, validation.Length(5, 256)),
	)
}

func NewUserCreatePayloadFromHttpRequest(req *http.Request) (*UserCreatePayload, error) {
	var item UserCreatePayload

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
