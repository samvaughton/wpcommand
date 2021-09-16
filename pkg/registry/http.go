package registry

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/samvaughton/wpcommand/v2/pkg/pipeline"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func GetHttpCallCommand(job types.CommandJob) pipeline.SiteCommand {
	return &pipeline.WrappedCommand{
		Name: CmdHttpCall,
		Wrapped: func(pipeline *pipeline.SiteCommandPipeline) (*types.CommandResult, error) {
			// perform a http call and then return the result
			req := http.Request{}

			u, err := url.Parse(job.Command.HttpUrl)

			if err != nil {
				return nil, errors.New(fmt.Sprintf("failed to parse the provided url: %s", job.Command.HttpUrl))
			}

			req.URL = u
			req.Method = job.Command.HttpMethod
			req.Body = ioutil.NopCloser(strings.NewReader(job.Command.HttpBody))
			req.Header = make(http.Header)

			var headers []map[string]string
			err = json.Unmarshal([]byte(job.Command.HttpHeaders), &headers)

			if err != nil {
				return nil, err
			}

			for _, headerPair := range headers {
				for headerKey, headerValue := range headerPair {
					req.Header.Set(headerKey, headerValue)
				}
			}

			client := &http.Client{}

			resp, err := client.Do(&req)

			if err != nil {
				return nil, errors.New(fmt.Sprintf("something went wrong performing the request: %s", err))
			}

			reqString, _ := httputil.DumpRequest(&req, true)
			respString, _ := httputil.DumpResponse(resp, true)

			cmdResult := &types.CommandResult{Command: CmdHttpCall, Output: fmt.Sprintf("%v", resp.StatusCode), Data: map[string]string{
				"Request":  string(reqString),
				"Response": string(respString),
			}}

			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				return cmdResult, errors.New(fmt.Sprintf("got non 2xx status code: %v", resp.StatusCode))
			}

			return cmdResult, nil
		},
	}
}
