package queue

import (
	"trading_engine/trading_engine"
)

// Dispatcher contains a pool of worker channels to process jobs
type Dispatcher interface {
	Run()
	Stop()
}

type dispatcher struct {
	MaxWorkers int
	WorkerPool chan chan Job
	Workers    []*Worker
	Engine     *trading_engine.TradingEngine
}

// NewDispatcher creates a new dispatcher
func NewDispatcher(engine *trading_engine.TradingEngine, maxWorkers int) Dispatcher {
	pool := make(chan chan Job, maxWorkers)
	workers := make([]*Worker, 0, maxWorkers)
	return &dispatcher{
		WorkerPool: pool,
		MaxWorkers: maxWorkers,
		Engine:     engine,
		Workers:    workers,
	}
}

// Run the dispatcher
func (d *dispatcher) Run() {
	// starting n number of workers
	for i := 0; i < d.MaxWorkers; i++ {
		worker := NewWorker(d.Engine, d.WorkerPool)
		worker.Start()
		d.Workers = append(d.Workers, &worker)
	}

	go d.dispatch()
}

func (d *dispatcher) dispatch() {
	for {
		select {
		case job := <-JobQueue:
			// a job request has been received
			go func(job Job) {
				// try to obtain a worker job channel that is available.
				// this will block until a worker is idle
				jobChannel := <-d.WorkerPool
				// dispatch the job to the worker job channel
				jobChannel <- job
			}(job)
		}
	}
}

// Stop the dispatcher
func (d *dispatcher) Stop() {
	for i := 0; i < d.MaxWorkers; i++ {
		d.Workers[i].Stop()
	}
}
