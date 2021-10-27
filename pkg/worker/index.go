package worker

import (
	"github.com/samvaughton/wpcommand/v2/pkg/pipeline"
	log "github.com/sirupsen/logrus"
	"sync/atomic"
	"time"
)

type CommandRunnerPool struct {
	MaxWorkers      int
	activeWorkers   int32
	workerCheck     *time.Ticker
	RecoveryHandler func(currentWorkerCount int32)
	queuedTaskChan  chan pipeline.SiteCommandPipeline
}

func NewCommandRunnerPool(maxWorkers int) *CommandRunnerPool {
	return &CommandRunnerPool{
		MaxWorkers:     maxWorkers,
		queuedTaskChan: make(chan pipeline.SiteCommandPipeline),
	}
}

func (wp *CommandRunnerPool) Start() {
	for i := 0; i < wp.MaxWorkers; i++ {
		go wp.startWorker()
	}

	wp.workerCheck = time.NewTicker(time.Second * 5)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-wp.workerCheck.C:
				if wp.GetActiveWorkers() < int32(wp.MaxWorkers) {
					// less active workers than there should be, lets add one at a time
					go wp.startWorker()
				}
			}
		}
	}()
}

func (wp *CommandRunnerPool) startWorker() {
	atomic.AddInt32(&wp.activeWorkers, 1) // increment worker count

	// handle panics
	defer func() {
		if err := recover(); err != nil {
			log.WithField("panic", "true").Error(err)
			newCount := atomic.AddInt32(&wp.activeWorkers, -1) // reduce active workers
			if wp.RecoveryHandler != nil {
				wp.RecoveryHandler(newCount)
			}
		}
	}()

	for item := range wp.queuedTaskChan {
		item.Run()
	}
}

func (wp *CommandRunnerPool) GetActiveWorkers() int32 {
	return atomic.LoadInt32(&wp.activeWorkers)
}

func (wp *CommandRunnerPool) AddTask(scp pipeline.SiteCommandPipeline) {
	wp.queuedTaskChan <- scp
}
