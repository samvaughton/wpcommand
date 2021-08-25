package types

import (
	"encoding/json"
	"gopkg.in/guregu/null.v3"
	"io/ioutil"
	"net/http"
	"time"
)

const CommandJobStatusCreated = "CREATED"
const CommandJobStatusPending = "PENDING"
const CommandJobStatusRunning = "RUNNING"
const CommandJobStatusSuccess = "SUCCESS"
const CommandJobStatusFailure = "FAILURE"

type ApiCreateCommandJobRequest struct {
	CommandId int64
	Selector  string
}

func NewApiCreateCommandJobRequest(req *http.Request) (*ApiCreateCommandJobRequest, error) {
	var data ApiCreateCommandJobRequest

	bytes, err := ioutil.ReadAll(req.Body)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &data)

	if err != nil {
		return nil, err
	}

	return &data, nil
}

type ApiCreateCommandJobResponse struct {
	Request   ApiCreateCommandJobRequest
	JobStatus string
	Sites     []*Site
	Jobs      []*CommandJob
}

type CommandJob struct {
	Id          int64 `bun:"id,pk"`
	RunByUserId null.Int
	RunByUser   *User `bun:"rel:belongs-to"`
	CommandId   int64
	Command     *Command `bun:"rel:belongs-to"`
	SiteId      int64
	Site        *Site `bun:"rel:belongs-to"`
	Uuid        string
	Key         string
	Status      string
	Description string
	CreatedAt   time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}
