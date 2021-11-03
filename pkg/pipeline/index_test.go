package pipeline

import (
	"github.com/pkg/errors"
	"github.com/samvaughton/wpcommand/v2/pkg/execution"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
	"testing"
)

func TestIsolatedCommands(t *testing.T) {
	exampleCommand := "example bash test command"

	testSite := &types.Site{TestMode: true}

	initDebugExecutor := func(site *types.Site) *execution.DebugCommandExecutor {
		exec, err := execution.NewCommandExecutor(testSite, &types.Config{})
		assert.Nil(t, err)
		debugExec, valid := exec.(*execution.DebugCommandExecutor)
		assert.True(t, valid)

		return debugExec
	}

	t.Run("TestSimpleCommand", func(t *testing.T) {
		exec := initDebugExecutor(testSite)

		cmds := []SiteCommand{
			&SimpleCommand{Args: []string{exampleCommand}},
		}

		pipeline := &SiteCommandPipeline{Name: "test pipeline", Site: testSite, Executor: exec, Options: ExecuteOptions{}, Commands: cmds}
		pipeline.Run()

		assert.Contains(t, exec.CommandLog, exampleCommand)
		exec.ClearLog()
	})

	t.Run("TestSimplePipelineCommand", func(t *testing.T) {
		exec := initDebugExecutor(testSite)

		cmds := []SiteCommand{
			&SimplePipelineCommand{
				Commands: []SiteCommand{
					&SimpleCommand{Args: []string{exampleCommand}},
				},
			},
		}

		pipeline := &SiteCommandPipeline{Name: "test pipeline", Site: testSite, Executor: exec, Options: ExecuteOptions{}, Commands: cmds}
		pipeline.Run()

		assert.Contains(t, exec.CommandLog, exampleCommand)
		exec.ClearLog()
	})

	t.Run("TestWrappedCommand", func(t *testing.T) {
		exec := initDebugExecutor(testSite)

		cmds := []SiteCommand{
			&WrappedCommand{Wrapped: func(pipeline *SiteCommandPipeline) (*types.CommandResult, error) {
				result, err := pipeline.Executor.ExecuteCommand([]string{exampleCommand})

				assert.Nil(t, err)

				return &types.CommandResult{Output: result.Output}, nil
			}},
		}

		pipeline := &SiteCommandPipeline{Name: "test pipeline", Site: testSite, Executor: exec, Options: ExecuteOptions{}, Commands: cmds}
		pipeline.Run()

		assert.Contains(t, exec.CommandLog, exampleCommand)
		exec.ClearLog()
	})

	t.Run("TestDynamicCommand", func(t *testing.T) {
		exec := initDebugExecutor(testSite)

		cmds := []SiteCommand{
			&DynamicArgsCommand{GetArgs: func(pipeline *SiteCommandPipeline) ([]string, error) {
				return []string{exampleCommand}, nil
			}},
		}

		pipeline := &SiteCommandPipeline{Name: "test pipeline", Site: testSite, Executor: exec, Options: ExecuteOptions{}, Commands: cmds}
		pipeline.Run()

		assert.Contains(t, exec.CommandLog, exampleCommand)
		exec.ClearLog()
	})
}

func TestPipelineFunctions(t *testing.T) {
	// Mock the function
	oGetPodBySite := execution.GetPodBySite
	execution.GetPodBySite = func(labelSelector string, namespace string, baseSelector string, k8Config *rest.Config) (*v1.Pod, error) {
		return &v1.Pod{}, nil
	}

	// After all tests, restore
	t.Cleanup(func() {
		execution.GetPodBySite = oGetPodBySite
	})

	exampleCommand := "example bash test command"

	testSite := &types.Site{TestMode: true}

	initDebugExecutor := func(site *types.Site) *execution.DebugCommandExecutor {
		exec, err := execution.NewCommandExecutor(testSite, &types.Config{})
		assert.Nil(t, err)
		debugExec, valid := exec.(*execution.DebugCommandExecutor)
		assert.True(t, valid)

		debugExec.MockOutput = "mock output"

		return debugExec
	}

	t.Run("TestHooksWithSuccessfulCommand", func(t *testing.T) {
		exec := initDebugExecutor(testSite)

		cmds := []SiteCommand{
			&SimpleCommand{Args: []string{exampleCommand}},
		}

		started := 0
		preAlways := 0
		skipped := 0
		postAlways := 0
		postSuccess := 0
		postError := 0
		finished := 0

		pipeline := &SiteCommandPipeline{
			Name:     "test pipeline",
			Site:     testSite,
			Executor: exec,
			Options:  ExecuteOptions{},
			Commands: cmds,
			Hooks: CommandHooks{
				Started: []func(){func() {
					started++
				}},
				PreAlways: []func(c SiteCommand){
					func(c SiteCommand) {
						preAlways++
					},
				},
				Skipped: []func(c SiteCommand, result *types.CommandResult, err error){
					func(c SiteCommand, result *types.CommandResult, err error) {
						skipped++
					},
				},
				PostAlways: []func(c SiteCommand, result *types.CommandResult, err error){
					func(c SiteCommand, result *types.CommandResult, err error) {
						postAlways++
					},
				},
				PostSuccess: []func(c SiteCommand, result *types.CommandResult, err error){
					func(c SiteCommand, result *types.CommandResult, err error) {
						postSuccess++
					},
				},
				PostError: []func(c SiteCommand, result *types.CommandResult, err error){
					func(c SiteCommand, result *types.CommandResult, err error) {
						postError++
					},
				},
				Finished: []func([]error){
					func(errors []error) {
						finished++
					},
				},
			},
		}

		pipeline.Run()

		assert.Contains(t, exec.CommandLog, exampleCommand)
		exec.ClearLog()

		assert.Equal(t, 1, started)
		assert.Equal(t, 1, preAlways)
		assert.Equal(t, 0, skipped)
		assert.Equal(t, 1, postAlways)
		assert.Equal(t, 0, postError)
		assert.Equal(t, 1, postSuccess)
		assert.Equal(t, 1, finished)
	})

	t.Run("TestErroneousHooks", func(t *testing.T) {
		exec := initDebugExecutor(testSite)

		cmds := []SiteCommand{
			&SimpleCommand{Args: []string{exampleCommand}},
			&WrappedCommand{Wrapped: func(pipeline *SiteCommandPipeline) (*types.CommandResult, error) {
				// also test previous output is returned
				prevOutput := pipeline.PreviousResult.Output

				assert.Equal(t, "mock output", prevOutput)

				return nil, errors.New("failure on purpose")
			}},
		}

		started := 0
		preAlways := 0
		skipped := 0
		postAlways := 0
		postSuccess := 0
		postError := 0
		finished := 0

		pipeline := &SiteCommandPipeline{
			Name:     "test pipeline",
			Site:     testSite,
			Executor: exec,
			Options:  ExecuteOptions{},
			Commands: cmds,
			Hooks: CommandHooks{
				Started: []func(){func() {
					started++
				}},
				PreAlways: []func(c SiteCommand){
					func(c SiteCommand) {
						preAlways++
					},
				},
				Skipped: []func(c SiteCommand, result *types.CommandResult, err error){
					func(c SiteCommand, result *types.CommandResult, err error) {
						skipped++
					},
				},
				PostAlways: []func(c SiteCommand, result *types.CommandResult, err error){
					func(c SiteCommand, result *types.CommandResult, err error) {
						postAlways++
					},
				},
				PostSuccess: []func(c SiteCommand, result *types.CommandResult, err error){
					func(c SiteCommand, result *types.CommandResult, err error) {
						postSuccess++
					},
				},
				PostError: []func(c SiteCommand, result *types.CommandResult, err error){
					func(c SiteCommand, result *types.CommandResult, err error) {
						postError++
					},
				},
				Finished: []func([]error){
					func(errors []error) {
						finished++
					},
				},
			},
		}

		pipeline.Run()

		assert.Contains(t, exec.CommandLog, exampleCommand)
		assert.Len(t, exec.CommandLog, 1)

		exec.ClearLog()

		assert.Equal(t, 1, started)
		assert.Equal(t, 2, preAlways)
		assert.Equal(t, 0, skipped)
		assert.Equal(t, 1, postError)
		assert.Equal(t, 2, postAlways)
		assert.Equal(t, 1, postSuccess)
		assert.Equal(t, 1, finished)
	})
}
