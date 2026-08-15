package core

import "context"

// Task represents a background executable job.
type Task interface {
	// ID returns a unique identifier for this background job.
	ID() string

	// Run contains the execution routine.
	Run(ctx context.Context) error
}

// Scheduler handles queuing and lifecycle tracking of asynchronous worker tasks.
type Scheduler interface {
	Component

	// Schedule registers and queues a task for immediate or delayed execution.
	Schedule(ctx context.Context, task Task) error

	// Cancel stops execution of a running task by ID.
	Cancel(ctx context.Context, taskID string) error

	// ListTasks returns active jobs running under the Scheduler.
	ListTasks(ctx context.Context) ([]Task, error)
}
