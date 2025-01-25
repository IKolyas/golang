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

	workerCh := make(chan func() error)

	stopCh := make(chan struct{})

	var result error

	wg := sync.WaitGroup{}

	mu := sync.Mutex{}

	for i := 0; i < n; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done() // Завершаем worker

			for {
				select {
				case <-stopCh:

					return

				case worker, ok := <-workerCh:

					if !ok {
						return
					}

					if err := worker(); err != nil {
						mu.Lock()

						m--

						mu.Unlock()
					}
				}
			}
		}()
	}

	for _, task := range tasks {
		mu.Lock()

		if m <= 0 {
			result = ErrErrorsLimitExceeded

			mu.Unlock()

			close(stopCh)

			break
		}

		mu.Unlock()

		workerCh <- task
	}

	close(workerCh) // Закрываем канал после отправки всех задач

	wg.Wait()

	return result
}
