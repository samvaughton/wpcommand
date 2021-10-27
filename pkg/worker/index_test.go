package worker

import (
	"errors"
	"github.com/samvaughton/wpcommand/v2/pkg/execution"
	"github.com/samvaughton/wpcommand/v2/pkg/pipeline"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	"github.com/stretchr/testify/assert"
	"sync/atomic"
	"testing"
	"time"
)

func TestProcessingAsExpected(t *testing.T) {
	t.Run("TestTasksAreProcessed", func(t *testing.T) {

		taskProcessed := int32(0)

		commandRunnerPool := NewCommandRunnerPool(3)

		commands := []pipeline.SiteCommand{
			&pipeline.WrappedCommand{
				Name: "Test Command",
				Wrapped: func(pipeline *pipeline.SiteCommandPipeline) (*types.CommandResult, error) {
					atomic.AddInt32(&taskProcessed, int32(1))

					return nil, nil
				},
			},
		}

		// Pass in empty stubs as no information should be pulled from these values for this test

		scp := pipeline.NewSiteCommandPipeline(
			"Test",
			&types.Site{},
			&types.Config{},
			commands,
			&execution.DebugCommandExecutor{Site: &types.Site{}, MockOutput: "Output Run"},
		)

		commandRunnerPool.Start()

		// add 10 tasks across 3 workers
		for i := 0; i < 10; i++ {
			commandRunnerPool.AddTask(scp)
		}

		assert.Eventually(t, func() bool {
			return taskProcessed == 10
		}, time.Second*10, time.Second*1)
	})

	t.Run("TestPanicIsHandled", func(t *testing.T) {
		taskProcessed := int32(0)

		commandRunnerPool := NewCommandRunnerPool(3)

		commands := []pipeline.SiteCommand{
			&pipeline.WrappedCommand{
				Name: "Test Command",
				Wrapped: func(pipeline *pipeline.SiteCommandPipeline) (*types.CommandResult, error) {
					atomic.AddInt32(&taskProcessed, int32(1))
					return nil, nil
				},
			},
		}

		// Pass in empty stubs as no information should be pulled from these values for this test

		scp := pipeline.NewSiteCommandPipeline(
			"Test",
			&types.Site{},
			&types.Config{},
			commands,
			&execution.DebugCommandExecutor{Site: &types.Site{}, MockOutput: "Output Run"},
		)

		commandRunnerPool.Start()

		// add 10 tasks across 3 workers
		for i := 0; i < 10; i++ {
			commandRunnerPool.AddTask(scp)
		}

		failingScp := scp
		failingScp.Commands = []pipeline.SiteCommand{
			&pipeline.WrappedCommand{
				Name: "Test Panic Command",
				Wrapped: func(pipeline *pipeline.SiteCommandPipeline) (*types.CommandResult, error) {
					panic(errors.New("test panic handling"))
				},
			},
		}

		commandRunnerPool.RecoveryHandler = func(currentWorkerCount int32) {
			assert.Equal(t, currentWorkerCount, int32(2)) // should be 2
		}

		commandRunnerPool.AddTask(failingScp)

		// we now need to check if the pool spools up a new worker
		assert.Eventually(t, func() bool {
			return commandRunnerPool.GetActiveWorkers() == int32(3)
		}, time.Second*10, time.Second*1)

		// now fire off another set of tasks to ensure that the new worker functions as intended
		for i := 0; i < 10; i++ {
			commandRunnerPool.AddTask(scp)
		}

		assert.Eventually(t, func() bool {
			return taskProcessed == 20 // two loops of 10
		}, time.Second*10, time.Second*1)
	})
}
