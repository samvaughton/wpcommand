package types

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/uptrace/bun"
	"io/ioutil"
	"net/http"
	"time"
)

type Site struct {
	Id            int64 `bun:"id,pk"`
	AccountId     int64
	Account       *Account        `bun:"rel:belongs-to" json:"-"`
	BlueprintSets []*BlueprintSet `bun:"m2m:sites_blueprint_sets" json:"-"`
	Uuid          string
	Key           string
	Description   string
	LabelSelector string
	Namespace     string
	SiteEmail     string
	SiteUsername  string
	SitePassword  string `casbin:"sitadwade,read_special"`
	SiteConfig    string
	Enabled       bool
	TestMode      bool
	CreatedAt     time.Time
	UpdatedAt     bun.NullTime
}

func NewSiteFromHttpRequest(req *http.Request) (*Site, error) {
	var site Site

	bytes, err := ioutil.ReadAll(req.Body)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &site)

	if err != nil {
		return nil, err
	}

	return &site, nil
}

type UpdateSitePayload struct {
	Description   string
	LabelSelector string
	Namespace     string
	Enabled       bool
	SiteConfig    string
	SiteUsername  string
	SiteEmail     string
	SitePassword  string
}

func (p UpdateSitePayload) HydrateSite(site *Site) {
	site.Description = p.Description
	site.Namespace = p.Namespace
	site.LabelSelector = p.LabelSelector
	site.Enabled = p.Enabled
	site.SiteConfig = p.SiteConfig
	site.SiteUsername = p.SiteUsername
	site.SiteEmail = p.SiteEmail
	site.SitePassword = p.SitePassword
}

func (p UpdateSitePayload) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Description, validation.Required),
		validation.Field(&p.LabelSelector, validation.Required),
		validation.Field(&p.Namespace, validation.Required),
		validation.Field(&p.SiteConfig, is.JSON),
		validation.Field(&p.SiteEmail, is.Email),
	)
}

func NewUpdateSitePayloadFromHttpRequest(req *http.Request) (*UpdateSitePayload, error) {
	var item UpdateSitePayload

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
