package types

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type DockerApiTagErrorResponse struct {
	Errinfo struct {
		Namespace  string `json:"namespace"`
		Repository string `json:"repository"`
		Tag        string `json:"tag"`
	} `json:"errinfo"`
	Message string `json:"message"`
}

type DockerApiTagResponse struct {
	Creator int         `json:"creator"`
	ID      int         `json:"id"`
	ImageID interface{} `json:"image_id"`
	Images  []struct {
		Architecture string      `json:"architecture"`
		Features     string      `json:"features"`
		Variant      interface{} `json:"variant"`
		Digest       string      `json:"digest"`
		Os           string      `json:"os"`
		OsFeatures   string      `json:"os_features"`
		OsVersion    interface{} `json:"os_version"`
		Size         int         `json:"size"`
		Status       string      `json:"status"`
		LastPulled   time.Time   `json:"last_pulled"`
		LastPushed   time.Time   `json:"last_pushed"`
	} `json:"images"`
	LastUpdated         time.Time `json:"last_updated"`
	LastUpdater         int       `json:"last_updater"`
	LastUpdaterUsername string    `json:"last_updater_username"`
	Name                string    `json:"name"`
	Repository          int       `json:"repository"`
	FullSize            int       `json:"full_size"`
	V2                  bool      `json:"v2"`
	TagStatus           string    `json:"tag_status"`
	TagLastPulled       time.Time `json:"tag_last_pulled"`
	TagLastPushed       time.Time `json:"tag_last_pushed"`
}

func NewDockerApiTagErrorFromResponse(resp *http.Response) (*DockerApiTagErrorResponse, error) {
	var item DockerApiTagErrorResponse

	bytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &item)

	if err != nil {
		return nil, err
	}

	return &item, nil
}

func NewDockerApiTagFromResponse(resp *http.Response) (*DockerApiTagResponse, error) {
	var item DockerApiTagResponse

	bytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &item)

	if err != nil {
		return nil, err
	}

	return &item, nil
}
