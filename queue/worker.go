package queue

import (
	"trading_engine/config"
	"trading_engine/trading_engine"
)

var (
	// MaxWorkers represents maximum number of workers to start
	MaxWorkers = config.Config.GetInt("MAX_WORKERS")
	// MaxQueue contains the maximum number of elements in the queue
	MaxQueue = config.Config.GetInt("MAX_QUEUE")
)

// Job structure to process
type Job struct {
	Order trading_engine.Order
}

// Worker contains information about the worker pool and job channel
type Worker struct {
	Engine     *trading_engine.TradingEngine
	WorkerPool chan chan Job
	JobChannel chan Job
	quit       chan bool
}

// JobQueue is a queue of channels
var JobQueue chan Job = make(chan Job, MaxQueue)

// NewWorker creates a new worker for the given worker pool
func NewWorker(engine *trading_engine.TradingEngine, workerPool chan chan Job) Worker {
	return Worker{
		Engine:     engine,
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan bool),
	}
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w Worker) Start() {
	go func() {
		for {
			// register the current worker into the worker queue.
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				// we have received a work request.
				w.Engine.Process(job.Order)

			case <-w.quit:
				// we have received a signal to stop
				return
			}
		}
	}()
}

// Stop signals the worker to stop listening for work requests.
func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
