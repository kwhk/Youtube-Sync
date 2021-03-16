package workerPool

import (
	"sync"
)

type Worker struct {
	ID int
	jobChan chan *Job
}

func NewWorker(jobChan chan *Job, ID int) *Worker {
	return &Worker{
		ID: ID,
		jobChan: jobChan,
	}
}

func (w *Worker) Run(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for job := range w.jobChan {
			job.Run()
		}
	}()
}