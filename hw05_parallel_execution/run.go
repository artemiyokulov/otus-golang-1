package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	tasksCount := int32(len(tasks))
	errorBudget := int32(m)
	recievedTasks := int32(0)
	neededWorkers := min(n, int(tasksCount))

	fifo := make(chan Task, tasksCount)

	waitGroup := sync.WaitGroup{}

	waitGroup.Add(neededWorkers)

	for i := 0; i < neededWorkers; i++ {
		go func() {
			defer waitGroup.Done()
			for {
				f, ok := <-fifo
				if !ok || errorBudget <= 0 || recievedTasks > tasksCount {
					return
				}
				if result := f(); result != nil {
					atomic.AddInt32(&errorBudget, -1)
				}
			}
		}()
	}

	for _, t := range tasks {
		fifo <- t
	}
	close(fifo)

	waitGroup.Wait()
	if errorBudget <= 0 {
		return ErrErrorsLimitExceeded
	}

	return nil
}
