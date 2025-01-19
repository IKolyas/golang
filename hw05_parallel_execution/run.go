package hw05parallelexecution

import (
	"context"
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	taskChan := make(chan func() error)
	errorChan := make(chan error)
	var wg sync.WaitGroup
	mu := &sync.Mutex{}

	ctx, cancel := context.WithCancel(context.Background())

	// Запускаем n горутин для обработки задач
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskChan {
				if err := task(); err != nil {
					mu.Lock()
					select {
					case errorChan <- err:
					case <-ctx.Done():
						return
					}
					mu.Unlock()
				}
			}
		}()
	}

	errorCount := 0
	go func() {
		for {
			select {
			case <-errorChan:
				mu.Lock()
				errorCount++
				if errorCount >= m {
					cancel()
					return
				}
				mu.Unlock()
			case <-ctx.Done():
				return
			}
		}
	}()

	for _, task := range tasks {
		taskChan <- task
	}

	// Завершаем работу канала тасков
	close(taskChan)

	wg.Wait()

	// После завершения работы канала тасков закрываем канал ошибок
	close(errorChan)

	if errorCount >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}
