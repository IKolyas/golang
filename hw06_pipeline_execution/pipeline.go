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

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := make(chan interface{})

	wg := sync.WaitGroup{}

	for _, stage := range stages {
		wg.Add(1)

		in = stage(in)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-done:
					return
				case value, ok := <-in:
					if !ok {
						return
					}
					out <- value
				case <-time.After(time.Second * 2):
					return
				}
			}
		}()
	}

	go func() {
		defer func() {
			for range in {
				<-in
			}
		}()
		wg.Wait()
		close(out)
	}()

	return out
}
