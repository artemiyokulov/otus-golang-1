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

func isTerminated(done In) bool {
	if done != nil {
		if _, ok := <-done; !ok {
			return true
		}
	}

	return false
}

func execChain(elem interface{}, done In, stages ...Stage) (result interface{}, ok bool) {
	ch := make(Bi, 1)
	ch <- elem
	close(ch)
	if len(stages) < 1 {
		return elem, true
	}
	processedValue, ok := <-(stages[0](ch))
	if !ok {
		log.Fatal("corrupted stage")
	}

	if isTerminated(done) {
		return nil, false
	}

	if len(stages) == 1 {
		return processedValue, true
	}
	return execChain(processedValue, done, stages[1:]...)
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
			if isTerminated(done) {
				wg.Add(-counter)
				return
			}
			elem, ok := <-in
			if !ok {
				return
			}
			wg.Add(1)
			counter++
			go func(elem interface{}, id int) {
				defer wg.Done()
				if isTerminated(done) {
					return
				}
				if result, ok := execChain(elem, done, stages...); ok {
					task := Task{
						id:     id,
						result: result,
					}
					outCh <- task
				} else {
					return
				}
			}(elem, counter)
		}
	}()
	return outCh
}
