package workerPool

// Task encapsulates a work item that should go in a work
// pool.
type Job struct {
	execute func()
}

func NewJob(f func()) *Job {
	return &Job{
		execute: f,
	}
}

func (job *Job) Run() {
	job.execute()
}