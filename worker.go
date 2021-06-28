package monitor

type WorkTask func(interface{}) interface{}

type WorkerPool struct {
	WorkerCount int
	Results     chan interface{}

	jobs     chan interface{}
	JobCount int
	work     WorkTask
}

func NewWorkerPool(workerCount int, jobs []interface{}, work WorkTask) *WorkerPool {
	// we shouldn't create more workers than there are jobs
	pool := WorkerPool{JobCount: len(jobs)}
	if workerCount > pool.JobCount {
		workerCount = pool.JobCount
	}

	pool.jobs = make(chan interface{}, pool.JobCount)
	pool.WorkerCount = workerCount
	pool.Results = make(chan interface{})
	pool.work = work

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
