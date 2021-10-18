package preview

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	"net/http"
)

func GetDockerTag(repo string, tag string) (*types.DockerApiTagResponse, error) {
	url := fmt.Sprintf("https://hub.docker.com/v2/repositories/%s/tags/%s", repo, tag)

	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		errorResp, err := types.NewDockerApiTagErrorFromResponse(resp)

		if err != nil {
			return nil, err
		}

		if resp.StatusCode == 404 {
			return nil, nil // not found
		}

		return nil, errors.New(errorResp.Message)
	}

	tagResp, err := types.NewDockerApiTagFromResponse(resp)

	if err != nil {
		return nil, err
	}

	return tagResp, nil
}
