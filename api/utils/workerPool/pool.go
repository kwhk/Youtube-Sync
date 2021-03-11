package workerPool

import (
	"sync"
	"fmt"
)

type Pool struct {
	JobsChan   chan *Job
	// Number of workers working concurrently
	concurrency int
	wg          sync.WaitGroup
}


// NewPool initializes a new pool with the given jobs and
// at the given concurrency.
func NewPool(concurrency int) *Pool {
	return &Pool{
		concurrency: concurrency,
		JobsChan:   make(chan *Job, 100),
	}
}

// Run runs all work within the pool and blocks until it's
// finished.
func (p *Pool) Run() {
	for i := 1; i <= p.concurrency; i++ {
		worker := NewWorker(p.JobsChan, i)
		worker.Run(&p.wg)
	}

	p.wg.Wait()
	close(p.JobsChan)
}

func (p *Pool) AddJob(f func()) {
	job := NewJob(f)
	p.JobsChan <- job
}