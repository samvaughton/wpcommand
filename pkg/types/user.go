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

type UserApiItem struct {
	Uuid       string
	Email      string
	FirstName  string
	LastName   string
	Account    string
	Roles      []string
	SuperAdmin bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (u *User) ConvertToApiItem() *UserApiItem {
	accountName := ""
	roles := []string{}

	if len(u.UserAccounts) == 1 {
		accountName = u.UserAccounts[0].Account.Name
		roles = u.UserAccounts[0].Roles
	}

	return &UserApiItem{
		Uuid:       u.Uuid,
		Email:      u.Email,
		FirstName:  u.FirstName,
		LastName:   u.LastName,
		Account:    accountName,
		Roles:      roles,
		SuperAdmin: u.SuperAdmin,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
}

type UserCreatePayload struct {
	Email     string
	FirstName string
	LastName  string
	Password  string
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

type AccountUserCreatePayload struct {
	Email     string
	FirstName string
	LastName  string
	Password  string
}

func (p AccountUserCreatePayload) Validate(uniqueCheck validation.Rule) error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Email, validation.Required, is.Email, uniqueCheck),
		validation.Field(&p.FirstName, validation.Required),
		validation.Field(&p.LastName, validation.Required),
		validation.Field(&p.Password, validation.Required, validation.Length(5, 256)),
	)
}

func NewAccountUserCreatePayloadFromHttpRequest(req *http.Request) (*AccountUserCreatePayload, error) {
	var item AccountUserCreatePayload

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

type AccountUserUpdatePayload struct {
	Email     string
	FirstName string
	LastName  string
	Password  string
}

func (p AccountUserUpdatePayload) Validate(uniqueCheck validation.Rule) error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Email, validation.Required, is.Email, uniqueCheck),
		validation.Field(&p.FirstName, validation.Required),
		validation.Field(&p.LastName, validation.Required),
		validation.Field(&p.Password, validation.Length(5, 256)),
	)
}

func NewAccountUserUpdatePayloadFromHttpRequest(req *http.Request) (*AccountUserUpdatePayload, error) {
	var item AccountUserUpdatePayload

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
