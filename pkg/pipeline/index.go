package pipeline

import (
	"github.com/samvaughton/wpcommand/v2/pkg/execution"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
)

type SiteCommand interface {
	getArgs(pipeline *SiteCommandPipeline) ([]string, error)
	GetName() string
	Execute(pipeline *SiteCommandPipeline) (*types.CommandResult, error)
}

type CommandHooks struct {
	Started     []func()                                                      // Is called when the Run() method is called
	PreAlways   []func(c SiteCommand)                                         // Is always called for each command at the start
	Skipped     []func(c SiteCommand, result *types.CommandResult, err error) // Is called when a command is skipped
	PostAlways  []func(c SiteCommand, result *types.CommandResult, err error) // Is always called after the command has run
	PostSuccess []func(c SiteCommand, result *types.CommandResult, err error) // Called after a successful command run
	PostError   []func(c SiteCommand, result *types.CommandResult, err error) // Called after an erroneous command run
	Finished    []func(errors []error)                                        // Called at the end of the Run() method
}

type SiteCommandPipeline struct {
	Name           string
	Site           *types.Site
	Config         *types.Config
	Hooks          CommandHooks
	Commands       []SiteCommand
	Options        ExecuteOptions
	Executor       execution.CommandExecutor
	PreviousResult *types.CommandResult
	Results        []*types.CommandResult
	Errors         []error
}

func NewSiteCommandPipeline(name string, site *types.Site, config *types.Config, commands []SiteCommand, executor execution.CommandExecutor) SiteCommandPipeline {
	return SiteCommandPipeline{
		Name:     name,
		Site:     site,
		Config:   config,
		Commands: commands,
		Executor: executor,
	}
}

type ExecuteOptions struct {
}

func (p *SiteCommandPipeline) Run() {
	for _, hook := range p.Hooks.Started {
		hook()
	}

	err := iterateAndExecute(p, p.Commands)

	if err != nil {
		p.Errors = append(p.Errors, err)
	}

	for _, hook := range p.Hooks.Finished {
		hook(p.Errors)
	}
}

func (p *SiteCommandPipeline) getArgs(pipeline *SiteCommandPipeline) ([]string, error) {
	panic("getArgs() on command pipeline should never be called")
}

func iterateAndExecute(pipeline *SiteCommandPipeline, commands []SiteCommand) error {
	for _, command := range commands {
		for _, hook := range pipeline.Hooks.PreAlways {
			hook(command)
		}

		result, err := command.Execute(pipeline)

		if err != nil {
			pipeline.Errors = append(pipeline.Errors, err)
		}

		preCheckError, isPreCheck := err.(*PreCheckFailedError)

		if isPreCheck {
			for _, hook := range pipeline.Hooks.Skipped {
				hook(command, result, preCheckError.Err)
			}
		} else {
			if err != nil {
				for _, hook := range pipeline.Hooks.PostError {
					hook(command, result, err)
				}
			} else {
				for _, hook := range pipeline.Hooks.PostSuccess {
					hook(command, result, nil)
				}
			}

			for _, hook := range pipeline.Hooks.PostAlways {
				hook(command, result, err)
			}
		}

		if err != nil {
			return err // any error we exit the command execution
		}

		pipeline.PreviousResult = result
		pipeline.Results = append(pipeline.Results, result)
	}

	return nil
}

func (p *SiteCommandPipeline) doExecute(c SiteCommand) (*types.CommandResult, error) {
	args, err := c.getArgs(p)

	if len(args) == 0 || err != nil {
		return nil, err
	}

	result, err := p.Executor.ExecuteCommand(args)

	return result, err
}
