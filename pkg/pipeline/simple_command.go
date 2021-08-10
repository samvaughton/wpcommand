package pipeline

import (
	"github.com/samvaughton/wpcommand/v2/pkg/types"
)

type SimpleCommand struct {
	SiteCommand
	Args []string
	Name string
}

func (c *SimpleCommand) getArgs(pipeline *SiteCommandPipeline) ([]string, error) {
	return c.Args, nil
}

func (c *SimpleCommand) GetName() string {
	return c.Name
}

func (c *SimpleCommand) Execute(pipeline *SiteCommandPipeline) (*types.CommandResult, error) {
	return pipeline.doExecute(c)
}
