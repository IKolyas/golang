package hw06pipelineexecution

import (
	"sync"
	"time"
)

type (
	In = <-chan interface{}

	Out = In

	Bi = chan interface{}
)

type Stage func(in In) (out Out)

func resAwait(in In, done In, out Bi) {
	for {
		select {
		case <-done:
			return
		case value, ok := <-in:
			if !ok {
				return
			}
			out <- value
		// Сливаем пайплайн по истечении отведённого времени
		case <-time.After(time.Second):
			return
		}
	}
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := make(chan interface{})

	wg := sync.WaitGroup{}

	for _, stage := range stages {
		in = stage(in)
	}

	wg.Add(1)

	go func() {
		defer wg.Done()
		resAwait(in, done, out)
	}()

	go func() {
		defer func() {
			for range in {
				// Сливаем зависшие вычисления
				// Ругается линтер
				<-in
			}
		}()

		wg.Wait()
		close(out)
	}()

	return out
}
