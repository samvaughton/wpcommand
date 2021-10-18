package types

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v3"
	"io/ioutil"
	"net/http"
	"time"
)

const CommandTypeHttpCall = "HTTP_CALL"
const CommandTypeWpBuiltIn = "WP_BUILT_IN"
const CommandTypePreviewBuild = "PREVIEW_BUILD"

type CommandTypePreviewBuildConfig struct {
	BuildPreviewRef string
}

func NewCommandTypePreviewBuildConfigFromConfig(cmdConfig string) (*CommandTypePreviewBuildConfig, error) {
	var item CommandTypePreviewBuildConfig

	err := json.Unmarshal([]byte(cmdConfig), &item)

	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (c CommandTypePreviewBuildConfig) ToString() (string, error) {
	cfgJson, err := json.Marshal(c)

	if err != nil {
		return "", err
	}

	return string(cfgJson), nil
}

type Command struct {
	Id          int64 `bun:"id,pk"`
	AccountId   null.Int
	Account     *Account `bun:"rel:belongs-to" json:"-"`
	SiteId      null.Int
	Site        *Site `bun:"rel:belongs-to" json:"-"`
	Public      bool
	Uuid        string
	Key         string
	Type        string
	Description string
	HttpMethod  string
	HttpUrl     string
	HttpHeaders string
	HttpBody    string
	Config      string
	CreatedAt   time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}

func (c *Command) IsDefault() bool {
	return c.AccountId.Valid == false && c.SiteId.Valid == false
}

type CreateCommandPayload struct {
	Type            string
	Description     string
	HttpMethod      string
	HttpUrl         string
	HttpHeaders     string
	HttpBody        string
	BuildPreviewRef string
	Public          bool
}

func (p CreateCommandPayload) HydrateCommand(command *Command) {
	command.Type = p.Type
	command.Uuid = uuid.New().String()
	command.Description = p.Description
	command.HttpMethod = p.HttpMethod
	command.HttpUrl = p.HttpUrl
	command.HttpHeaders = p.HttpHeaders
	command.HttpBody = p.HttpBody
	command.Public = p.Public

	if command.HttpHeaders == "" {
		command.HttpHeaders = "{}"
	}

	cfgJson, err := json.Marshal(map[string]string{
		"BuildPreviewRef": p.BuildPreviewRef,
	})

	if err != nil {
		log.Error(err)
	} else {
		command.Config = string(cfgJson)
	}
}

func (p CreateCommandPayload) Validate() error {
	if p.Type == CommandTypeHttpCall {
		return validation.ValidateStruct(&p,
			validation.Field(&p.Type, validation.Required),
			validation.Field(&p.Description, validation.Required),
			validation.Field(&p.HttpUrl, validation.Required, is.URL),
			validation.Field(&p.HttpMethod, validation.Required),
			validation.Field(&p.HttpHeaders, is.JSON),
		)
	} else if p.Type == CommandTypePreviewBuild {
		return validation.ValidateStruct(&p,
			validation.Field(&p.Type, validation.Required),
			validation.Field(&p.Description, validation.Required),
			validation.Field(&p.BuildPreviewRef, validation.Required),
		)
	}

	return validation.ValidateStruct(&p,
		validation.Field(&p.Type, validation.Required),
		validation.Field(&p.Description, validation.Required),
	)
}

func NewCreateCommandPayloadFromHttpRequest(req *http.Request) (*CreateCommandPayload, error) {
	var item CreateCommandPayload

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

type UpdateCommandPayload struct {
	Type            string
	Description     string
	HttpMethod      string
	HttpUrl         string
	HttpHeaders     string
	HttpBody        string
	BuildPreviewRef string
	Public          bool
}

func (p UpdateCommandPayload) HydrateCommand(command *Command) {
	command.Type = p.Type
	command.Uuid = uuid.New().String()
	command.Description = p.Description
	command.HttpMethod = p.HttpMethod
	command.HttpUrl = p.HttpUrl
	command.HttpHeaders = p.HttpHeaders
	command.HttpBody = p.HttpBody
	command.Public = p.Public

	cfgJson, err := json.Marshal(map[string]string{
		"BuildPreviewRef": p.BuildPreviewRef,
	})

	if err != nil {
		log.Error(err)
	} else {
		command.Config = string(cfgJson)
	}

	if command.HttpHeaders == "" {
		command.HttpHeaders = "{}"
	}

	if command.Config == "" {
		command.Config = "{}"
	}
}

func (p UpdateCommandPayload) Validate() error {
	if p.Type == CommandTypeHttpCall {
		return validation.ValidateStruct(&p,
			validation.Field(&p.Type, validation.Required),
			validation.Field(&p.Description, validation.Required),
			validation.Field(&p.HttpUrl, validation.Required, is.URL),
			validation.Field(&p.HttpMethod, validation.Required),
			validation.Field(&p.HttpHeaders, is.JSON),
		)
	} else if p.Type == CommandTypePreviewBuild {
		return validation.ValidateStruct(&p,
			validation.Field(&p.Type, validation.Required),
			validation.Field(&p.Description, validation.Required),
			validation.Field(&p.BuildPreviewRef, validation.Required),
		)
	}

	return validation.ValidateStruct(&p,
		validation.Field(&p.Type, validation.Required),
		validation.Field(&p.Description, validation.Required),
	)
}

func NewUpdateCommandPayloadFromHttpRequest(req *http.Request) (*UpdateCommandPayload, error) {
	var item UpdateCommandPayload

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
