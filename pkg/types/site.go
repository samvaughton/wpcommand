package types

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
	"io/ioutil"
	"net/http"
	"time"
)

type Site struct {
	ApiSiteCore
	ApiSiteCredentials
}

type ApiSiteCore struct {
	Id                 int64 `bun:"id,pk"`
	AccountId          int64
	Account            *Account        `bun:"rel:belongs-to" json:"-"`
	BlueprintSets      []*BlueprintSet `bun:"m2m:sites_blueprint_sets" json:"-"`
	Uuid               string
	Key                string
	Description        string
	LabelSelector      string
	Namespace          string
	SiteEmail          string
	SiteUsername       string
	SitePassword       string
	SiteConfig         string
	WpCachedData       string
	WpDomain           string
	DockerRegistryName string

	Enabled   bool
	TestMode  bool
	CreatedAt time.Time
	UpdatedAt bun.NullTime
}

func NewApiSiteCoreFromSite(site Site) *ApiSiteCore {
	return &ApiSiteCore{
		Id:                 site.Id,
		AccountId:          site.AccountId,
		Account:            site.Account,
		BlueprintSets:      site.BlueprintSets,
		Uuid:               site.Uuid,
		Key:                site.Key,
		Description:        site.Description,
		LabelSelector:      site.LabelSelector,
		Namespace:          site.Namespace,
		SiteEmail:          site.SiteEmail,
		SiteUsername:       site.SiteUsername,
		SitePassword:       site.SitePassword,
		SiteConfig:         site.SiteConfig,
		WpCachedData:       site.WpCachedData,
		WpDomain:           site.WpDomain,
		DockerRegistryName: site.DockerRegistryName,

		Enabled:   site.Enabled,
		TestMode:  site.TestMode,
		CreatedAt: site.CreatedAt,
		UpdatedAt: site.UpdatedAt,
	}
}

type ApiSiteCredentials struct {
	AccessToken string
}

func NewApiSiteCredentialsFromSite(site Site) *ApiSiteCredentials {
	return &ApiSiteCredentials{
		AccessToken: site.AccessToken,
	}
}

func (s *Site) GetWpCachedData() (WpCachedData, error) {
	var data WpCachedData

	err := json.Unmarshal([]byte(s.WpCachedData), &data)

	if err != nil {
		log.Errorf("failed decoding wp cached data: %s", err)

		return data, err
	}

	return data, nil
}

func (s *Site) SetWpCachedData(data *WpCachedData) error {
	bytes, err := json.Marshal(data)

	if err != nil {
		return err
	}

	s.WpCachedData = string(bytes)

	return nil
}

type WpCachedData struct {
	UserList []WpUser `json:"UserList"`
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
	Description        string
	LabelSelector      string
	Namespace          string
	Enabled            bool
	SiteConfig         string
	SiteUsername       string
	SiteEmail          string
	SitePassword       string
	WpDomain           string
	DockerRegistryName string
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
	site.WpDomain = p.WpDomain
	site.DockerRegistryName = p.DockerRegistryName
}

func (p UpdateSitePayload) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Description, validation.Required),
		validation.Field(&p.LabelSelector, validation.Required),
		validation.Field(&p.Namespace, validation.Required),
		validation.Field(&p.WpDomain, validation.Required),
		validation.Field(&p.DockerRegistryName, validation.Required),
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
