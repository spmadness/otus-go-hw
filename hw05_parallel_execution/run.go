package hw05parallelexecution

import (
	"errors"
	"sync"
)

var (
	ErrErrorsLimitExceeded  = errors.New("errors limit exceeded")
	ErrWorkersValueTooSmall = errors.New("workers value must be greater than 0")
	ErrTasksSliceIsEmpty    = errors.New("tasks slice is empty")
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errCnt int
	var err error

	if m <= 0 {
		return ErrErrorsLimitExceeded
	}
	if n <= 0 {
		return ErrWorkersValueTooSmall
	}
	if len(tasks) == 0 {
		return ErrTasksSliceIsEmpty
	}

	taskCh := make(chan Task, len(tasks))

	for _, t := range tasks {
		taskCh <- t
	}
	close(taskCh)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for v := range taskCh {
				res := v()
				if res == nil {
					continue
				}

				mu.Lock()
				if errCnt == m {
					mu.Unlock()
					return
				}
				errCnt++
				if errCnt < m {
					mu.Unlock()
					continue
				}

				err = ErrErrorsLimitExceeded
				mu.Unlock()
				return
			}
		}()
	}

	wg.Wait()

	if err != nil {
		return err
	}
	return nil
}
