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

	if tasksCount == 0 {
		return nil
	}

	errorBudget := int32(m)

	fifo := make(chan Task, tasksCount)

	waitGroup := sync.WaitGroup{}

	waitGroup.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			defer waitGroup.Done()
			for f := range fifo {
				if result := f(); result != nil {
					atomic.AddInt32(&errorBudget, -1)
				}
				if atomic.LoadInt32(&errorBudget) <= 0 {
					return
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
