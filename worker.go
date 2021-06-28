package monitor

import (
	"runtime"

	"github.com/thomasrea0113/gpu-price-monitor/domain"
)

type WorkTask func(domain.PriceCheckJob) domain.PriceCheckResponse

type WorkerPool struct {
	Results chan domain.PriceCheckResponse

	jobs     chan domain.PriceCheckJob
	JobCount int
	work     WorkTask
}

func NewWorkerPool(jobs []domain.PriceCheckJob, work WorkTask) *WorkerPool {
	pool := WorkerPool{JobCount: len(jobs)}
	pool.Results = make(chan domain.PriceCheckResponse)
	pool.work = work

	pool.jobs = make(chan domain.PriceCheckJob, pool.JobCount)
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
	for c := 0; c < runtime.NumCPU(); c++ {
		go pool.doWork()
	}
}
