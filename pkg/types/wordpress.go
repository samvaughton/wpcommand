package types

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"io/ioutil"
	"net/http"
)

type WpUser struct {
	ID             int    `json:"ID"`
	Roles          string `json:"roles"`
	DisplayName    string `json:"display_name"`
	UserLogin      string `json:"user_login"`
	UserEmail      string `json:"user_email"`
	UserStatus     string `json:"user_status"`
	UserRegistered string `json:"user_registered"`
}

type ApiWpUser struct {
	Id       int
	Username string
	Email    string
	Password string
	Role     string
}

func NewApiWpUserListFromWpUserList(users []WpUser) []*ApiWpUser {
	var items = make([]*ApiWpUser, 0)

	for _, u := range users {
		items = append(items, NewApiWpUserFromWpUser(u))
	}

	return items
}

func NewApiWpUserFromWpUser(user WpUser) *ApiWpUser {
	return &ApiWpUser{
		Id:       user.ID,
		Role:     user.Roles,
		Email:    user.UserEmail,
		Username: user.UserLogin,
	}
}

type UpdateWpUserPayload struct {
	Email    string
	Password string
	Role     string
}

func (p UpdateWpUserPayload) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Email, validation.Required, is.Email),
		validation.Field(&p.Password, validation.Length(6, 256)),
		validation.Field(&p.Role, validation.Required, validation.In("owner")),
	)
}

func NewUpdateWpUserPayloadFromHttpRequest(req *http.Request) (*UpdateWpUserPayload, error) {
	var item UpdateWpUserPayload

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

type CreateWpUserPayload struct {
	Username string
	Email    string
	Password string
	Role     string
}

func (p CreateWpUserPayload) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Email, validation.Required, is.Email),
		validation.Field(&p.Username, validation.Required),
		validation.Field(&p.Password, validation.Required, validation.Length(6, 256)),
		validation.Field(&p.Role, validation.Required, validation.In("owner")),
	)
}

func NewCreateWpUserPayloadFromHttpRequest(req *http.Request) (*CreateWpUserPayload, error) {
	var item CreateWpUserPayload

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

type WpPlugin struct {
	Name    string `json:"name" yaml:"name"`
	Status  string `json:"status" yaml:"status"`
	Update  string `json:"update" yaml:"update"`
	Version string `json:"version" yaml:"version"`
	Url     string `json:"url" yaml:"url"`
}

type WpTheme struct {
	Name    string `json:"name" yaml:"name"`
	Status  string `json:"status" yaml:"status"`
	Update  string `json:"update" yaml:"update"`
	Version string `json:"version" yaml:"version"`
	Url     string `json:"url" yaml:"url"`
}

type WpPost struct {
	Id     int    `json:"ID"`
	Title  string `json:"post_title"`
	Name   string `json:"post_name"`
	Date   string `json:"post_date"`
	Status string `json:"post_status"`
}
