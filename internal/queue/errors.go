package queue

import "errors"

var (
	// ErrQueueFull is returned when the job queue is full
	ErrQueueFull = errors.New("job queue is full")
)