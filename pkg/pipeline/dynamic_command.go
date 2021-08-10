package pipeline

import (
	"github.com/samvaughton/wpcommand/v2/pkg/types"
)

type DynamicArgsCommand struct {
	Name string
	SiteCommand
	GetArgs func(pipeline *SiteCommandPipeline) ([]string, error)
}

func (c *DynamicArgsCommand) getArgs(pipeline *SiteCommandPipeline) ([]string, error) {
	return c.GetArgs(pipeline)
}

func (c *DynamicArgsCommand) GetName() string {
	return c.Name
}

func (c *DynamicArgsCommand) Execute(pipeline *SiteCommandPipeline) (*types.CommandResult, error) {
	return pipeline.doExecute(c)
}
