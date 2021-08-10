package pipeline

import (
	"github.com/samvaughton/wpcommand/v2/pkg/types"
)

type SimplePipelineCommand struct {
	SiteCommand
	Name               string
	ErrorIsSuccessful  bool
	RunPreCheckCommand SiteCommand
	Commands           []SiteCommand
}

type PreCheckFailedError struct {
	Err error
}

func (u *PreCheckFailedError) Error() string {
	return u.Err.Error()
}

func (c *SimplePipelineCommand) GetName() string {
	return c.Name
}

func (c *SimplePipelineCommand) Execute(pipeline *SiteCommandPipeline) (*types.CommandResult, error) {
	if c.RunPreCheckCommand != nil {
		result, err := c.RunPreCheckCommand.Execute(pipeline)

		if (c.ErrorIsSuccessful && err != nil) || (!c.ErrorIsSuccessful && err == nil) {
			iterateAndExecute(pipeline, c.Commands)

			return nil, nil
		}

		return result, &PreCheckFailedError{Err: err}
	} else {
		iterateAndExecute(pipeline, c.Commands)
	}

	return nil, nil
}
