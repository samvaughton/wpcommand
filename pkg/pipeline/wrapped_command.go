package pipeline

import "github.com/samvaughton/wpcommand/v2/pkg/types"

type WrappedCommand struct {
	SiteCommand
	Name    string
	Wrapped func(pipeline *SiteCommandPipeline) (*types.CommandResult, error)
}

func (c *WrappedCommand) getArgs(pipeline *SiteCommandPipeline) ([]string, error) {
	return []string{}, nil
}

func (c *WrappedCommand) GetName() string {
	return c.Name
}

func (c *WrappedCommand) Execute(pipeline *SiteCommandPipeline) (*types.CommandResult, error) {
	result, err := c.Wrapped(pipeline)

	if err != nil {
		return result, err
	}

	return result, nil
}
