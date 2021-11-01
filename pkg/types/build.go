package types

import "github.com/google/go-github/v39/github"

type BuildRequest struct {
	Id                 string
	RepoOwner          string
	RepoName           string
	RepoRef            string
	Workflow           string
	DockerRegistryName string
	WordpressDomain    string
	IsPreviewBuild     bool
}

type BuildResult struct {
	BuildRequest      BuildRequest
	GithubWorkflowJob github.WorkflowJob
	BuildUrl          string
}
