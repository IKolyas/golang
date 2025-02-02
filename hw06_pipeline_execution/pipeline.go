package hw06pipelineexecution

import (
	"sync"
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
		select {
		case <-done:
			return nil
		default:
			in = stage(in)
		}
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			case val, ok := <-in:
				if !ok {
					return
				}
				select {
				case <-done:
					return
				case out <- val:
				}
			}
		}
	}()

	go func() {
		defer func() {
			close(out)
			for range in {
				<-in
			}
		}()
		wg.Wait()
	}()

	return out
}
