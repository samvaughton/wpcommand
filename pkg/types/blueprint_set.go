package types

import (
	"encoding/json"
	"github.com/uptrace/bun"
	"io/ioutil"
	"net/http"
	"time"
)

type BlueprintSet struct {
	Id        int64 `bun:"id,pk"`
	AccountId int64
	Account   *Account `bun:"rel:belongs-to" json:"-"`
	Sites     []*Site  `bun:"m2m:sites_blueprint_sets"`
	Uuid      string
	Name      string
	Enabled   bool
	CreatedAt time.Time
	UpdatedAt bun.NullTime
}

type CreateBlueprintSetPayload struct {
	Name string
}

func NewCreateBlueprintSetPayloadFromHttpRequest(req *http.Request) (*CreateBlueprintSetPayload, error) {
	var item CreateBlueprintSetPayload

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
