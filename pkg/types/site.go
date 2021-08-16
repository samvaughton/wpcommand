package types

import (
	"encoding/json"
	"github.com/uptrace/bun"
	"io/ioutil"
	"net/http"
	"time"
)

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
	SitePassword  string
	SiteConfig    string
	Enabled       bool
	TestMode      bool
	CreatedAt     time.Time
	UpdatedAt     bun.NullTime
}
