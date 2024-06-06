package hw06pipelineexecution

import (
	"log"
	"sync"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

type Task struct {
	id     int
	result interface{}
}

func execChain(elem interface{}, done In, stages ...Stage) interface{} {
	ch := make(Bi, 1)
	ch <- elem
	close(ch)
	if len(stages) < 1 {
		return elem
	}
	processedValue, ok := <-(stages[0](ch))
	if !ok {
		log.Fatal("corrupted stage")
	}
	select {
	case <-done:
		return nil
	default:
		if len(stages) == 1 {
			return processedValue
		}
		return execChain(processedValue, done, stages[1:]...)
	}
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	outCh := make(Bi)
	go func() {
		wg := sync.WaitGroup{}
		counter := 0
		defer func() {
			wg.Wait()
			close(outCh)
		}()
		for {
			select {
			case <-done:
				wg.Add(-counter)
				return
			case elem, ok := <-in:
				if !ok {
					return
				}
				wg.Add(1)
				counter++
				go func(elem interface{}, id int) {
					defer wg.Done()
					task := Task{
						id:     id,
						result: execChain(elem, done, stages...),
					}
					select {
					case <-done:
						return
					default:
						outCh <- task
					}
				}(elem, counter)
			}
		}

	}()
	return outCh
}
