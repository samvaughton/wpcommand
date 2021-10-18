package types

import "github.com/google/go-github/v39/github"

type BuildPreviewRequest struct {
	Id                 string
	RepoOwner          string
	RepoName           string
	RepoRef            string
	Workflow           string
	DockerRegistryName string
	WordpressDomain    string
}

type BuildPreviewResult struct {
	BuildPreviewRequest BuildPreviewRequest
	GithubWorkflowJob   github.WorkflowJob
	BuildPreviewUrl     string
}
