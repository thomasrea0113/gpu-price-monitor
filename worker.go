package monitor

import (
	"sync"
)

type WorkTask func(interface{}) interface{}

type WorkerPool struct {
	WorkerCount int
	Results     chan interface{}

	jobs     chan interface{}
	JobCount int
	work     WorkTask
}

func NewWorkerPool(workerCount int, jobs []interface{}, work WorkTask) *WorkerPool {
	pool := WorkerPool{JobCount: len(jobs)}
	pool.WorkerCount = workerCount
	pool.Results = make(chan interface{})
	pool.work = work

	pool.jobs = make(chan interface{}, pool.JobCount)
	for _, job := range jobs {
		pool.jobs <- job
	}
	close(pool.jobs)

	return &pool
}

func (pool WorkerPool) doWork() {
	for j := range pool.jobs {
		pool.Results <- pool.work(j)
	}
}

func (pool WorkerPool) Start() {
	for c := 0; c < pool.WorkerCount; c++ {
		go pool.doWork()
	}
}

func (pool WorkerPool) StartGroup() *sync.WaitGroup {
	var wg sync.WaitGroup
	for c := 0; c < pool.WorkerCount; c++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			pool.doWork()
		}()
	}

	return &wg
}
