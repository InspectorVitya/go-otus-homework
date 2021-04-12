package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
var ErrInvalidNumberWorkers = errors.New("invalid number workers")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if err := checkParamsInput(n, m); err != nil {
		return err
	}

	var ch = make(chan Task, len(tasks))
	wg := sync.WaitGroup{}
	maxErrs := int32(m)
	var errCount int32
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			stop := false
			for val := range ch {
				if stop {
					break
				}
				if err := val(); err != nil {
					stop = atomic.AddInt32(&errCount, 1) >= maxErrs
				} else {
					stop = atomic.LoadInt32(&errCount) >= maxErrs
				}
			}
		}()
	}
	for _, task := range tasks {
		ch <- task
	}
	close(ch)
	wg.Wait()
	if errCount >= maxErrs {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func checkParamsInput(n, m int) error {
	switch {
	case n <= 0:
		return ErrInvalidNumberWorkers
	case m <= 0:
		return ErrErrorsLimitExceeded
	default:
		return nil
	}
}
