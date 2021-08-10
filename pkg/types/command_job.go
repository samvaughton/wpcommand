package types

import (
	"encoding/json"
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
	Command  string
	Selector string
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
	Sites     []Site
}

type CommandJob struct {
	Id          int64 `bun:"id,pk"`
	SiteId      int64
	Site        *Site `bun:"rel:belongs-to"`
	Uuid        string
	Key         string
	Status      string
	Description string
	CreatedAt   time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}
