package worker

import (
	"github.com/samvaughton/wpcommand/v2/pkg/pipeline"
)

type CommandRunnerPool struct {
	MaxWorkers     int
	queuedTaskChan chan pipeline.SiteCommandPipeline
}

func NewCommandRunnerPool(maxWorkers int) *CommandRunnerPool {
	return &CommandRunnerPool{
		MaxWorkers:     maxWorkers,
		queuedTaskChan: make(chan pipeline.SiteCommandPipeline),
	}
}

func (wp *CommandRunnerPool) Start() {
	for i := 0; i < wp.MaxWorkers; i++ {
		go func(workerId int) {
			for item := range wp.queuedTaskChan {
				item.Run()
			}
		}(i + 1)
	}
}

func (wp *CommandRunnerPool) AddTask(scp pipeline.SiteCommandPipeline) {
	wp.queuedTaskChan <- scp
}
