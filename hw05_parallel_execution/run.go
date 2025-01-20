package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) error {
	if m < 0 {
		m = 0
	}

	taskCh := make(chan func() error)
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}

	var res error

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskCh {
				if err := task(); err != nil {
					mu.Lock()
					m--
					if m <= 0 {
						res = ErrErrorsLimitExceeded
						mu.Unlock()
						return
					}
					mu.Unlock()
				}
			}
		}()
	}

	for _, task := range tasks {
		taskCh <- task
	}
	close(taskCh) // Закрываем канал после отправки всех задач

	wg.Wait()

	return res
}
