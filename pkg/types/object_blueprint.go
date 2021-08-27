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

const ObjectBlueprintTypePlugin = "PLUGIN"
const ObjectBlueprintTypeTheme = "THEME"

type ObjectBlueprint struct {
	Id         int64 `bun:"id,pk"`
	RevisionId int64

	BlueprintSetId int64
	BlueprintSet   *BlueprintSet `bun:"rel:belongs-to"`

	Uuid     string
	SetOrder int
	Type     string
	Name     string
	Enabled  bool

	Version   string
	ExactName string

	OriginalObjectUrl   string
	VersionedObjectUrl  string
	VersionedObjectHash string

	CreatedAt time.Time
	UpdatedAt bun.NullTime
}

func NewCreateObjectBlueprintPayloadFromHttpRequest(req *http.Request) (*CreateObjectBlueprintPayload, error) {
	var item CreateObjectBlueprintPayload

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

type CreateObjectBlueprintPayload struct {
	Type      string
	Name      string
	ExactName string
	Version   string
	Url       string
}

func (p CreateObjectBlueprintPayload) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Type, validation.Required, validation.In(ObjectBlueprintTypePlugin, ObjectBlueprintTypeTheme)),
		validation.Field(&p.Name, validation.Required),
		validation.Field(&p.ExactName, validation.Required),
		validation.Field(&p.Version, validation.Required, is.Semver),
		validation.Field(&p.Url, validation.Required, is.URL),
	)
}
