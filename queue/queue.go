package queue

import "github.com/Frankuccino/relay/job"

type Queue struct {
	jobs chan job.Job
}

func New(bufferSize int) *Queue {
	return &Queue{
		jobs: make(chan job.Job, bufferSize),
	}
}

func (q *Queue) Enqueue(j job.Job) {
	q.jobs <- j
}

func (q *Queue) Jobs() <-chan job.Job {
	return q.jobs
}

func (q *Queue) Len() int {
	return len(q.jobs)
}
