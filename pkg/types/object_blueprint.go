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

	ObjectBlueprintStorage []*ObjectBlueprintStorage `bun:"m2m:object_blueprint_storage_relations" json:"-"`

	Uuid     string
	SetOrder int
	Type     string
	Name     string

	Active bool

	Version   string
	ExactName string

	OriginalObjectUrl string

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

func NewUpdateObjectBlueprintPayloadFromHttpRequest(req *http.Request) (*UpdateObjectBlueprintPayload, error) {
	var item UpdateObjectBlueprintPayload

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

func NewUpdatedVersionObjectBlueprintPayloadFromHttpRequest(req *http.Request) (*UpdatedVersionObjectBlueprintPayload, error) {
	var item UpdatedVersionObjectBlueprintPayload

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

type UpdateObjectBlueprintPayload struct {
	Name string
}

type CreateObjectBlueprintPayload struct {
	Type      string
	Name      string
	ExactName string
	Version   string
	Url       string
}

type UpdatedVersionObjectBlueprintPayload struct {
	Version string
	Url     string
}

func (p UpdateObjectBlueprintPayload) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Name, validation.Required),
	)
}

func (p UpdatedVersionObjectBlueprintPayload) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Version, validation.Required, is.Semver),
	)
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
