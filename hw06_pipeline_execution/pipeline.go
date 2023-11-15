package hw06pipelineexecution

import (
	"sync"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := make(chan interface{})
	wg := sync.WaitGroup{}
	valChs := make([]chan interface{}, 0)

	if in == nil || stages == nil {
		close(out)
		return out
	}

	for v := range in {
		wg.Add(1)

		// для каждого входного значения создаем канал и передаем его в горутину, для того,
		// чтобы значения были на выходе в том же порядке, что и на входе.
		valCh := make(chan interface{}, 1)
		valChs = append(valChs, valCh)

		go processStages(v, valCh, done, stages, &wg)
	}

	// ждем обработки всех значений, перебираем каналы со значениями в исходном порядке,
	// освобождаем функцию ExecutePipeline
	go waitAndSend(valChs, out, &wg)

	return out
}

func processStages(v interface{}, valCh chan interface{}, done In, stages []Stage, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := range stages {
		if stages[i] == nil {
			continue
		}
		ch := make(chan interface{}, 1)
		ch <- v
		select {
		case <-done:
			close(ch)
			close(valCh)
			return
		case v = <-stages[i](ch):
			close(ch)
		}
	}
	valCh <- v
}

func waitAndSend(valChs []chan interface{}, out chan interface{}, wg *sync.WaitGroup) {
	wg.Wait()

	for _, ch := range valChs {
		if v, ok := <-ch; ok {
			out <- v
			close(ch)
		}
	}

	close(out)
}
