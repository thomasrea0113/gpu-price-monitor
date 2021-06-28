package monitor

import "github.com/thomasrea0113/gpu-price-monitor/domain"

type WorkTask func(domain.PriceCheckJob) domain.PriceCheckResponse

type WorkerPool struct {
	Results chan domain.PriceCheckResponse

	jobs     chan domain.PriceCheckJob
	jobCount int
	work     WorkTask
}

func NewWorkerPool(jobs []domain.PriceCheckJob, work WorkTask) *WorkerPool {
	pool := WorkerPool{jobCount: len(jobs)}
	pool.jobs = make(chan domain.PriceCheckJob, pool.jobCount)
	for _, job := range jobs {
		pool.jobs <- job
	}
	close(pool.jobs)

	pool.Results = make(chan domain.PriceCheckResponse)

	return &pool
}

func (pool WorkerPool) DoWork() {
	for j := range pool.jobs {
		// TODO resolve segfault after first worker starts
		pool.Results <- pool.work(j)
	}
}
