package pipeline

import (
	"fmt"
	"github.com/samvaughton/wpcommand/v2/pkg/execution"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
)

type SiteCommand interface {
	getArgs(pipeline *SiteCommandPipeline) ([]string, error)
	GetName() string
	Execute(pipeline *SiteCommandPipeline) (*types.CommandResult, error)
}

type CommandHooks struct {
	Started     []func()
	PreAlways   []func(c SiteCommand)
	Skipped     []func(c SiteCommand, result *types.CommandResult, err error)
	PostAlways  []func(c SiteCommand, result *types.CommandResult, err error)
	PostSuccess []func(c SiteCommand, result *types.CommandResult, err error)
	PostError   []func(c SiteCommand, result *types.CommandResult, err error)
	Finished    []func(errors []error)
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
	errors         []error
}

type ExecuteOptions struct {
}

func (p *SiteCommandPipeline) Run() {
	for _, hook := range p.Hooks.Started {
		hook()
	}

	err := iterateAndExecute(p, p.Commands)

	if err != nil {
		p.errors = append(p.errors, err)
	}

	for _, hook := range p.Hooks.Finished {
		hook(p.errors)
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
			fmt.Println("exec", err)
			pipeline.errors = append(pipeline.errors, err)
		}

		preCheckError, isPreCheck := err.(*PreCheckFailedError)

		if isPreCheck {
			for _, hook := range pipeline.Hooks.Skipped {
				hook(command, result, preCheckError.Err)
			}
		} else {
			if err != nil {
				fmt.Println("hook", err)
				for _, hook := range pipeline.Hooks.PostError {
					hook(command, result, err)
				}

				break
			} else {
				for _, hook := range pipeline.Hooks.PostSuccess {
					hook(command, result, nil)
				}
			}
		}

		if err != nil {
			fmt.Println("last check", err)
			return err // always break on an error
		}

		pipeline.PreviousResult = result
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
